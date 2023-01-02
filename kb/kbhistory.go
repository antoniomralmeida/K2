package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *KBHistory) Persist() error {
	ctx, collection := initializers.GetCollection("KBHistory")
	if h.Id.IsZero() {
		h.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, h)
		return err
	} else {
		_, err := collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: h.Id}}, h)
		return err
	}
}

func (h *KBHistory) ClearingHistory(history int) error {

	type PipeCount struct {
		Id    primitive.ObjectID `json:"_id"`
		Count int                `json:"count"`
	}

	Id := h.Attribute
	ctx, collection := initializers.GetCollection("KBHistory")
	for {
		matchStage := bson.D{{Key: "attribute_id", Value: Id}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$attribute_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
		ret, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
		initializers.Log(err, initializers.Error)
		if err != nil {
			return nil
		}
		results := []PipeCount{}
		err = ret.All(ctx, &results)
		initializers.Log(err, initializers.Error)
		if err != nil {
			return nil
		}
		if results[0].Count <= history {
			return nil
		}
		todel := KBHistory{}
		collection.FindOne(ctx, bson.D{{Key: "attribute_id", Value: Id}}, options.FindOne().SetSort(bson.D{{Key: "when", Value: 1}})).Decode(&todel)
		if !todel.Id.IsZero() {
			collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: todel.Id}})
		} else {
			return nil
		}
	}
}

func (h *KBHistory) FindLast(filter bson.D) error {
	ctx, collection := initializers.GetCollection("KBHistory")
	ret := collection.FindOne(ctx, filter, options.FindOne().SetSort(bson.D{{Key: "when", Value: -1}}))
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
