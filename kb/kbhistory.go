package kb

import (
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *KBHistory) Persist() error {
	if h.ID.IsZero() {
		err := mgm.Coll(h).Create(h)
		return err
	} else {
		db_doc := new(KBHistory)
		err := mgm.Coll(h).FindByID(h.ID, db_doc)
		if err != nil {
			return err
		}
		if h.UpdatedAt != db_doc.UpdatedAt {
			return errors.New("Old document!")
		}
		err = mgm.Coll(h).Update(h)
		return err
	}
}

func (h *KBHistory) ClearingHistory(history int) error {

	type PipeCount struct {
		Id    primitive.ObjectID `json:"_id"`
		Count int                `json:"count"`
	}

	Id := h.Attribute
	collection := mgm.Coll(h)
	for {
		matchStage := bson.D{{Key: "attribute_id", Value: Id}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$attribute_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
		ret, err := collection.Aggregate(mgm.Ctx(), mongo.Pipeline{matchStage, groupStage})
		initializers.Log(err, initializers.Error)
		if err != nil {
			return nil
		}
		results := []PipeCount{}
		err = ret.All(mgm.Ctx(), &results)
		initializers.Log(err, initializers.Error)
		if err != nil {
			return nil
		}
		if results[0].Count <= history {
			return nil
		}
		todel := KBHistory{}
		collection.FindOne(mgm.Ctx(), bson.D{{Key: "attribute_id", Value: Id}}, options.FindOne().SetSort(bson.D{{Key: "when", Value: 1}})).Decode(&todel)
		if !todel.ID.IsZero() {
			collection.DeleteOne(mgm.Ctx(), bson.D{{Key: "_id", Value: todel.ID}})
		} else {
			return nil
		}
	}
}

func (h *KBHistory) FindLast(filter bson.D) error {
	collection := mgm.Coll(h)
	ret := collection.FindOne(mgm.Ctx(), filter, options.FindOne().SetSort(bson.D{{Key: "when", Value: -1}}))
	if ret != nil {
		ret.Decode(h)
	}
	return nil
}

func (h *KBHistory) String() string {
	j, err := json.MarshalIndent(*h, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
