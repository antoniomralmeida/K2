package models

import (
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KBAlert struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Message          string               `bson:"message"`
	User             primitive.ObjectID   `bson:"user"`
	Views            []primitive.ObjectID `bson:"views"`
}

func getLastAlerts(user primitive.ObjectID) ([]KBAlert, error) {
	alerts := []KBAlert{}
	opts := options.Find().SetSort(bson.D{{"create_at", -1}}).SetLimit(5)

	err := mgm.Coll(new(KBAlert)).SimpleFind(alerts, bson.M{"$or": bson.A{
		bson.M{"user": user},
		bson.M{"user": ""},
	}}, opts)
	return alerts, err
}

func NewAlert(msg string, email string) error {
	alert := new(KBAlert)
	alert.Message = msg
	if email != "" {
		user := new(KBUser)
		err := user.FindOne(bson.D{{Key: "email", Value: email}})
		if err != nil {
			return err
		}
		alert.User = user.ID
	}
	err := alert.Persist()
	return err
}

func (obj *KBAlert) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (obj *KBAlert) Persist() error {
	return inits.Persist(obj)

}
