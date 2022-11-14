package kb

import (
	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func (w *KBWorkspace) Persist() error {
	ctx, collection := initializers.GetCollection("KBWorkspace")
	if w.Id.IsZero() {
		w.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, w)
		return err
	} else {
		_, err := collection.UpdateOne(ctx, bson.D{{Name: "_id", Value: w.Id}}, w)
		return err
	}
}

func FindAllWorkspaces(sort string) error {
	ctx, collection := initializers.GetCollection("KBWorkspace")
	idx := collection.Indexes()
	ret, err := idx.List(ctx)
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		_, err = idx.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"workspace": 1}, Options: options.Index().SetUnique(true)})
		initializers.Log(err, initializers.Fatal)
	}
	return nil
}
