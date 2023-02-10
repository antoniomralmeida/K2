package models

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KBObject struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string              `bson:"name"`
	Class            primitive.ObjectID  `bson:"class_id"`
	Attributes       []KBAttributeObject `bson:"attributes"`
	Bkclass          *KBClass            `bson:"-" json:"Class"`
	parsed           bool                `bson:"-"`
}

func ObjectFactory(class string, name string) *KBObject {
	p := FindClassByName(class, true)
	if p == nil {
		inits.Log("Class not found "+class, inits.Error)
		return nil
	}
	o := KBObject{Name: name, Class: p.ID, Bkclass: p}
	for _, x := range FindAttributes(p) {
		n := KBAttributeObject{Attribute: x.ID, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		//_kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	inits.Log(o.Persist(), inits.Fatal)
	//_kb.IdxObjects[name] = &o
	return &o
}

func FindObjectByName(name string) (ret *KBObject) {
	r := mgm.Coll(ret).FindOne(mgm.Ctx(), bson.D{{"name", name}})
	r.Decode(ret)
	return
}

func ObjectFacroryByClass(name string, class *KBClass) *KBObject {
	o := KBObject{Name: name, Class: class.ID, Bkclass: class}
	for _, x := range FindAttributes(class) {
		n := KBAttributeObject{Attribute: x.ID, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		//_kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	inits.Log(o.Persist(), inits.Fatal)
	//_kb.IdxObjects[name] = &o
	return &o
}

func (obj *KBObject) Persist() error {
	return inits.Persist(obj)

}

func (obj *KBObject) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (o *KBObject) String() string {
	j, err := json.MarshalIndent(*o, "", "\t")
	inits.Log(err, inits.Error)
	return string(j)
}

func (o *KBObject) Delete() error {

	mgm.Coll(o).Delete(o)

	// Restart KB
	KBStop()
	KBInit()
	return nil
}
