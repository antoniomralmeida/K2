package models

import (
	"time"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StackStatus byte

const (
	Pending StackStatus = iota
	Running
	Concluded
)

type KBStack struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	RuleID           primitive.ObjectID `bson:"rule_id"`
	Status           StackStatus        `bson:"status"`
}

func (obj *KBStack) Persist() error {
	return inits.Persist(obj)
}
func (obj *KBStack) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func StackFactory(id primitive.ObjectID) *KBStack {
	stack := new(KBStack)
	stack.RuleID = id
	inits.Log(stack.Persist(), inits.Fatal)
	return stack
}

func KBAddStack(rules []*KBRule) error {
	for _, r := range rules {
		StackFactory(r.ID)
	}
	return nil
}

func StacktoRun() (list []KBRule) {
	mgm.Coll(new(KBStack)).UpdateMany(mgm.Ctx(), bson.D{{Key: "status", Value: Pending}}, bson.D{{"status", Running}, {"update_at", time.Now().UTC()}})
	ret, err := mgm.Coll(new(KBStack)).Distinct(mgm.Ctx(), "rule_id", bson.D{{"status", Running}})
	inits.Log(err, inits.Error)
	oids := make([]primitive.ObjectID, len(ret))
	for _, id := range ret {
		oids = append(oids, id.(primitive.ObjectID))
	}
	opts := options.Find().SetSort(bson.D{{Key: "priority", Value: 1}, {Key: "lastexecution", Value: -1}})
	mgm.Coll(new(KBRule)).SimpleFind(list, bson.D{{Key: "$in", Value: oids}}, opts)
	return
}

func StackEndRun() error {
	_, err := mgm.Coll(new(KBStack)).UpdateMany(mgm.Ctx(), bson.D{{"status", Running}}, bson.D{{"status", Concluded}, {"update_at", time.Now().UTC()}})
	return err
}
