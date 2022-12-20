package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (o *KBObject) Persist() error {
	ctx, collection := initializers.GetCollection("KBObject")
	if o.Id.IsZero() {
		o.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, o)
		return err
	} else {
		_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: o.Id}}, o)
		return err
	}
}

func (o *KBObject) String() string {
	j, err := json.MarshalIndent(*o, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}

func (o *KBObject) Delete() error {

	ctx, collection := initializers.GetCollection("KBObject")
	collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: o.Id}})

	// Restart KB
	Stop()
	Init()
	return nil
}

func (o *KBObject) GetWorkspaces() (ret []*KBWorkspace) {
	for i := range GKB.Workspaces {
		for j := range GKB.Workspaces[i].Objects {
			if GKB.Workspaces[i].Objects[j].KBObject == o {
				ret = append(ret, &GKB.Workspaces[i])
			}
		}
	}
	return
}
