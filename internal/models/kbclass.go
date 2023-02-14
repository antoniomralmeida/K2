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

func KBClassFactory(name, icon, parent string) (class *KBClass, err error) {
	if parent != "" {
		parentClass := FindClassByName(parent, true)
		if parentClass == nil {
			inits.Log(lib.ClassNotFoundError, inits.Info)
			return nil, lib.ClassNotFoundError
		}
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
	ok, err := a.Valitate()
	inits.Log(err, inits.Error)
	if ok {
		obj.Attributes = append(obj.Attributes, a)
		err = obj.Persist()
		if err == nil {
			return &a, nil
		}
	}
	inits.Log(err, inits.Fatal)
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
	x := mgm.Coll(class).FindOne(mgm.Ctx(), p)
	if x != nil {
		x.Decode(class)
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
				KBRestart()
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

func FindAttributes(c *KBClass) []*KBAttribute {
	var ret []*KBAttribute
	if c != nil {
		if c.ParentClass != nil {
			ret = append(ret, FindAttributes(c.ParentClass)...)
		}
		for i := range c.Attributes {
			ret = append(ret, &c.Attributes[i])
		}
	}
	return ret
}

func FindAttribute(c *KBClass, name string) *KBAttribute {
	attrs := FindAttributes(c)
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
