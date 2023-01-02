package kb

import (
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (class *KBClass) Persist() error {
	ctx, collection := initializers.GetCollection("KBClass")
	if class.Id.IsZero() {
		class.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, class)
		return err
	} else {
		_, err := collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: class.Id}}, class)
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

func (class *KBClass) Delete(force bool) error {
	//TODO: Verificar se há classes filhas, se houver não exclui
	//TODO: Verificar se há objetos, se houver não exclui
	//TODO: com force, excluir todas as dependências antes
	ctx, collection := initializers.GetCollection("KBClass")
	collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: class.Id}})

	// Restart KB
	Stop()
	Init()
	return nil
}

func (class *KBClass) String() string {
	j, err := json.MarshalIndent(*class, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
