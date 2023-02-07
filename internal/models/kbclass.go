package models

import (
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KBClass struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string             `bson:"name"`
	Icon             string             `bson:"icon"`
	ParentID         primitive.ObjectID `bson:"parent_id,omitempty"`
	Parent           string             `bson:"-" json:"parent"`
	Attributes       []KBAttribute      `bson:"attributes"`
	ParentClass      *KBClass           `bson:"-"`
}

func ClassFactory(newclass_json string) *KBClass {
	class := KBClass{}
	err := json.Unmarshal([]byte(newclass_json), &class)
	if err != nil {
		inits.Log(err, inits.Info)
		return nil
	}
	if class.Parent != "" {
		p := _kb.FindClassByName(class.Parent, true)
		if p == nil {
			inits.Log("Class not found "+class.Parent, inits.Info)
			return nil
		}
		class.ParentID = p.ID
		class.ParentClass = p
	}
	for i := range class.Attributes {
		class.Attributes[i].ID = primitive.NewObjectID()
		for _, x := range class.Attributes[i].Sources {
			class.Attributes[i].SourcesID = append(class.Attributes[i].SourcesID, KBSourceStr[x])
		}
		class.Attributes[i].SimulationID = KBSimulationStr[class.Attributes[i].Simulation]
	}
	err = class.Persist()
	if err == nil {
		_kb.Classes = append(_kb.Classes, class)
		_kb.IdxClasses[class.ID] = &class
		return &class
	} else {
		inits.Log(err, inits.Error)
		return nil
	}
}

func (obj *KBClass) Persist() error {
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
