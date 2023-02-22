package models

import (
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

type KBObject struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string              `bson:"name" valid:"length(5|50),required"`
	Class            primitive.ObjectID  `bson:"class_id"`
	Attributes       []KBAttributeObject `bson:"attributes"`
	ClassPtr         *KBClass            `bson:"-"`
	Parsed           bool                `bson:"-"`
}

func (obj *KBObject) validate() (bool, error) {
	return govalidator.ValidateStruct(obj)
}

func ObjectFactory(name, className string) (*KBObject, error) {
	class := FindClassByName(className, true)
	if class == nil {
		inits.Log(lib.ClassNotFoundError, inits.Error)
		return nil, lib.ClassNotFoundError
	}
	return ObjectFactoryByClass(name, class)
}

func (obj *KBObject) validateIndex() error {
	cur, err := mgm.Coll(obj).Indexes().List(mgm.Ctx())
	inits.Log(err, inits.Error)
	var result []bson.M
	err = cur.All(mgm.Ctx(), &result)
	if len(result) == 1 {
		inits.CreateUniqueIndex(mgm.Coll(obj), "name")
	}
	return err
}

func FindObjectByName(name string) (ret *KBObject) {
	ret = nil
	cur := mgm.Coll(ret).FindOne(mgm.Ctx(), bson.D{{"name", name}})
	inits.Log(cur.Err(), inits.Error)
	if cur.Err() == nil {
		ret = new(KBObject)
		cur.Decode(ret)
	}
	return
}

func ObjectFactoryByClass(name string, class *KBClass) (*KBObject, error) {
	obj := KBObject{Name: name, Class: class.ID, ClassPtr: class}
	for _, x := range class.FindAttributes() {
		n := KBAttributeObject{ID: primitive.NewObjectID(), Attribute: x.ID, KbAttribute: x, KbObject: &obj}
		obj.Attributes = append(obj.Attributes, n)
	}
	ok, err := obj.validate()
	if !ok {
		inits.Log(err, inits.Error)
		return nil, err
	}
	err = obj.Persist()
	if mongo.IsDuplicateKeyError(err) {
		inits.Log(err, inits.Error)
	} else {
		inits.Log(err, inits.Fatal)
	}
	return &obj, err
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
	if kb_current != nil {
		// Restart KB
		kb_current.RestartFlag = true
	}
	return nil
}

func FindAllObjects(filter bson.M, sort string) (objs []KBObject, err error) {
	objs = []KBObject{}
	cursor, err := mgm.Coll(new(KBObject)).Find(mgm.Ctx(), filter, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	if err == nil {
		err = cursor.All(mgm.Ctx(), &objs)
	}
	return
}
