package kb

import (
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (obj *KBClass) Persist() error {
	return initializers.Persist(obj)

}

func (obj *KBClass) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (class *KBClass) FindOne(p bson.D) error {
	x := mgm.Coll(class).FindOne(mgm.Ctx(), p)
	if x != nil {
		x.Decode(class)
		return nil
	} else {
		return errors.New("Class not found!")
	}
}

func (class *KBClass) Delete(force bool) error {
	//TODO: Verificar se há classes filhas, se houver não exclui
	//TODO: Verificar se há objetos, se houver não exclui
	//TODO: com force, excluir todas as dependências antes
	mgm.Coll(class).Delete(class)

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
