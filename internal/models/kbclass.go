package models

import (
	"context"
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/asaskevich/govalidator"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KBClass struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string             `bson:"name" valid:"length(5|50),required"`
	Icon             string             `bson:"icon"`
	ParentID         primitive.ObjectID `bson:"parent_id,omitempty"`
	Parent           string             `bson:"-" json:"parent"`
	Attributes       []KBAttribute      `bson:"attributes"`
	ParentClass      *KBClass           `bson:"-"`
}

func (obj *KBClass) validateIndex() error {
	cur, err := mgm.Coll(obj).Indexes().List(mgm.Ctx())
	inits.Log(err, inits.Error)
	var result []bson.M
	err = cur.All(context.TODO(), &result)
	if len(result) == 1 {
		inits.CreateUniqueIndex(mgm.Coll(obj), "name")
	}
	return err
}

func (obj *KBClass) valitate() (bool, error) {
	return govalidator.ValidateStruct(obj)
}

func KBClassFactoryParent(name, icon string, parentClass *KBClass) (class *KBClass, err error) {
	if parentClass != nil {
		class = &KBClass{Name: name, Icon: icon, ParentClass: parentClass, ParentID: parentClass.ID}
	} else {
		class = &KBClass{Name: name, Icon: icon}
	}
	ok, err := class.valitate()
	inits.Log(err, inits.Error)
	if !ok {
		return nil, err
	}
	err = class.Persist()
	if err != nil {
		return nil, err
	}
	return class, nil
}

func KBClassFactory(name, icon, parent string) (class *KBClass, err error) {
	var parentClass *KBClass
	if parent != "" {
		parentClass = FindClassByName(parent, true)
		if parentClass == nil {
			inits.Log(lib.ClassNotFoundError, inits.Info)
			return nil, lib.ClassNotFoundError
		}
	} else {
		parentClass = nil
	}
	return KBClassFactoryParent(name, icon, parentClass)
}

func (obj *KBClass) AlterClassAddAttribute(name, atype, simulation string, options, sources []string, keephistory int, valitade int64) (attr *KBAttribute, err error) {
	a := KBAttribute{ID: primitive.NewObjectID(),
		Name:             name,
		AType:            KBattributeTypeStr(atype),
		Options:          options,
		Sources:          sources,
		SourcesID:        ToKBSources(sources),
		KeepHistory:      keephistory,
		ValidityInterval: valitade,
		Simulation:       simulation,
		SimulationID:     KBSimulationStr[simulation]}
	ok, err := a.validate()
	inits.Log(err, inits.Error)
	if ok {
		obj.Attributes = append(obj.Attributes, a)
		err = obj.Persist()
		if err == nil {
			return &a, nil
		}
		inits.Log(err, inits.Error)
	}
	return nil, err
}

func (obj *KBClass) Persist() error {
	obj.validateIndex()
	return inits.Persist(obj)
}

func (obj *KBClass) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (class *KBClass) FindOne(p bson.D) error {
	ret := mgm.Coll(class).FindOne(mgm.Ctx(), p)
	if ret.Err() == nil {
		ret.Decode(class)
		return nil
	} else {
		return lib.ClassNotFoundError
	}
}

func (class *KBClass) Delete() error {
	res := mgm.Coll(class).FindOne(mgm.Ctx(), bson.D{{"parente", class.ID}})
	if res.Err() == mongo.ErrNoDocuments {
		res = mgm.Coll(new(KBObject)).FindOne(mgm.Ctx(), bson.D{{"class", class.ID}})
		if res.Err() == mongo.ErrNoDocuments {
			err := mgm.Coll(class).Delete(class)
			if err == nil {
				restartKB()
			}
			return err
		}
	}
	return mongo.ErrMultipleIndexDrop
}

func (class *KBClass) String() string {
	j, err := json.MarshalIndent(*class, "", "\t")
	inits.Log(err, inits.Error)
	return string(j)
}

func (c *KBClass) FindAttributes() []*KBAttribute {
	var ret []*KBAttribute
	if c.ParentClass != nil {
		ret = append(ret, c.ParentClass.FindAttributes()...)
	}
	for i := range c.Attributes {
		ret = append(ret, &c.Attributes[i])
	}
	return ret
}

func (c *KBClass) FindAttribute(name string) *KBAttribute {
	attrs := c.FindAttributes()
	for i, x := range attrs {
		if x.Name == name {
			return attrs[i]
		}
	}
	return nil
}

func FindClassByName(nm string, mandatory bool) *KBClass {
	ret := new(KBClass)
	err := ret.FindOne(bson.D{{Key: "name", Value: nm}})
	if err != nil && mandatory {
		inits.Log(err, inits.Error)
		return nil
	}
	return ret
}

func FindAllClasses(sort string) (class *[]KBClass, err error) {
	class = new([]KBClass)
	cursor, err := mgm.Coll(new(KBClass)).Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	if err == nil {
		err = cursor.All(mgm.Ctx(), class)
	}
	return
}

func KBClassCopy(name string, copy *KBClass) (*KBClass, error) {
	if copy == nil {
		inits.Log(lib.InvalidClassError, inits.Error)
		return nil, lib.InvalidClassError
	}
	class := *copy
	class.ID = primitive.NilObjectID
	class.Name = name
	for i := range class.Attributes {
		class.Attributes[i].ID = primitive.NewObjectID()
	}
	err := class.Persist()
	if err == nil {
		_classes = append(_classes, class)
		return &class, nil
	} else {
		inits.Log(err, inits.Error)
		return nil, err
	}
}
