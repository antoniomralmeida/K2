package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"gopkg.in/mgo.v2/bson"
)

func (class *KBHistory) Persist() error {
	collection := initializers.GetDb().C("KBHistory")
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}

func (h *KBHistory) FindLast(filter bson.D) error {
	collection := initializers.GetDb().C("KBHistory")
	return collection.Find(filter).Sort("-when").One(&h)
}

func (h *KBHistory) String() string {
	j, err := json.MarshalIndent(*h, "", "\t")
	lib.LogFatal(err)
	return string(j)
}
