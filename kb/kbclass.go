package kb

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllClasses(sort string, cs *[]KBClass) error {
	ctx, collection := initializers.GetCollection("KBClass")
	idx := collection.Indexes()
	ret, err := idx.List(context.TODO())
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		_, err = idx.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)})
		initializers.Log(err, initializers.Fatal)
	}
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, cs)
	return err
}

func (class *KBClass) Persist() error {
	ctx, collection := initializers.GetCollection("KBClass")
	if class.Id.IsZero() {
		class.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, class)
		return err
	} else {
		_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: class.Id}}, class)
		return err
	}
}

func (class *KBClass) FindOne(p bson.D) error {
	ctx, collection := initializers.GetCollection("KBClass")
	x := collection.FindOne(ctx, p)
	if x != nil {
		x.Decode(class)
		return nil
	} else {
		return errors.New("Not found Class!")
	}
}

func (class *KBClass) Delete() error {
	//TODO: Restart KB
	return nil
}

func (class *KBClass) String() string {
	j, err := json.MarshalIndent(*class, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
