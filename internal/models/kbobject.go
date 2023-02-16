package models

import (
	"context"
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	for _, x := range p.FindAttributes() {
		n := KBAttributeObject{Attribute: x.ID, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		//_kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	inits.Log(o.Persist(), inits.Fatal)
	//_kb.IdxObjects[name] = &o
	return &o
}

func (obj *KBObject) validateIndex() error {
	cur, err := mgm.Coll(obj).Indexes().List(mgm.Ctx())
	inits.Log(err, inits.Error)
	var result []bson.M
	err = cur.All(context.TODO(), &result)
	if len(result) == 1 {
		inits.CreateUniqueIndex(mgm.Coll(obj), "name")
	}
	return err
}

func FindObjectByName(name string) (ret *KBObject) {
	r := mgm.Coll(ret).FindOne(mgm.Ctx(), bson.D{{"name", name}})
	r.Decode(ret)
	return
}

func ObjectFactoryByClass(name string, class *KBClass) *KBObject {
	o := KBObject{Name: name, Class: class.ID, Bkclass: class}
	for _, x := range class.FindAttributes() {
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
	stopKB()
	InitKB()
	return nil
}

func FindAllObjects(filter bson.M, sort string, objs *[]KBObject) error {
	cursor, err := mgm.Coll(new(KBObject)).Find(mgm.Ctx(), filter, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), objs)
	return err
}
