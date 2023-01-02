package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (w *KBWorkspace) Persist() error {
	ctx, collection := initializers.GetCollection("KBWorkspace")
	if w.Id.IsZero() {
		w.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, w)
		return err
	} else {
		_, err := collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: w.Id}}, w)
		return err
	}
}

func (w *KBWorkspace) String() string {
	j, err := json.MarshalIndent(*w, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
