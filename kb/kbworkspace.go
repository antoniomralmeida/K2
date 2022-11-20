package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (w *KBWorkspace) Persist() error {
	ctx, collection := initializers.GetCollection("KBWorkspace")
	if w.Id.IsZero() {
		w.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, w)
		return err
	} else {
		_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: w.Id}}, w)
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
	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, &GKB.Workspaces)
	return err
}

func (w *KBWorkspace) String() string {
	j, err := json.MarshalIndent(*w, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
