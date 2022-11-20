package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	//TODO: Reiniciar KB
	return nil
}

func FindAllObjects(filter bson.M, sort string, os *[]KBObject) error {
	ctx, collection := initializers.GetCollection("KBObject")
	idx := collection.Indexes()
	ret, err := idx.List(ctx)
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		_, err = idx.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)})
		initializers.Log(err, initializers.Fatal)
	}
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, os)
	return err
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
