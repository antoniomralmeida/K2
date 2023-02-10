package models

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/asaskevich/govalidator"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KBClass struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string             `bson:"name" valid:"length(5|50)"`
	Icon             string             `bson:"icon"`
	ParentID         primitive.ObjectID `bson:"parent_id,omitempty"`
	Parent           string             `bson:"-" json:"parent"`
	Attributes       []KBAttribute      `bson:"attributes"`
	ParentClass      *KBClass           `bson:"-"`
}

func (obj *KBClass) validateIndex() error {
	cur, err := mgm.Coll(obj).Indexes().List(mgm.Ctx())
	var result []bson.M
	err = cur.All(context.TODO(), &result)
	if len(result) == 1 {
		inits.CreateUniqueIndex(mgm.Coll(obj), "name")
	}
	return err
}

func (obj *KBClass) valitade() (bool, error) {
	return govalidator.ValidateStruct(obj)
}

func ClassFactory(name, icon, parent string) (class *KBClass) {
	if parent != "" {
		parentClass := FindClassByName(parent, true)
		if parentClass == nil {
			inits.Log("Class not found "+class.Parent, inits.Info)
			return nil
		}
		class = &KBClass{Name: name, Icon: icon, ParentClass: parentClass, ParentID: parentClass.ID}
	} else {
		class = &KBClass{Name: name, Icon: icon}
	}
	ok, err := class.valitade()
	inits.Log(err, inits.Error)
	if !ok {
		return nil
	}
	err = class.Persist()
	if err != nil {
		return nil
	}
	return class
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
		return errors.New("Class not found!")
	}
}

func (class *KBClass) Delete() error {
	res := mgm.Coll(class).FindOne(mgm.Ctx(), bson.D{{"parente", class.ID}})
	if res.Err() == mongo.ErrNoDocuments {
		res = mgm.Coll(new(KBObject)).FindOne(mgm.Ctx(), bson.D{{"class", class.ID}})
		if res.Err() == mongo.ErrNoDocuments {
			err := mgm.Coll(class).Delete(class)
			if err == nil {
				// Restart KB
				KBStop()
				KBInit()
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

func (c *KBClass) UpdateClass() {
	for i := range c.Attributes {
		if c.Attributes[i].ID.IsZero() {
			c.Attributes[i].ID = primitive.NewObjectID()
		}
	}
	inits.Log(c.Persist(), inits.Fatal)
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
