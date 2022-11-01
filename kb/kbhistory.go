package kb

import (
	"github.com/antoniomralmeida/k2/initializers"
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
