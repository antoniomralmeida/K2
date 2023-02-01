package models

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
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
