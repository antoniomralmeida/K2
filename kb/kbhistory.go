package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"gopkg.in/mgo.v2/bson"
)

func (h *KBHistory) Persist() error {
	collection := initializers.GetDb().C("KBHistory")
	if h.Id == "" {
		h.Id = bson.NewObjectId()
		return collection.Insert(h)
	} else {
		return collection.UpdateId(h.Id, h)
	}
}

func (h *KBHistory) ClearingHistory(history int) error {
	Id := h.Attribute
	collection := initializers.GetDb().C("KBHistory")
	for {
		n, err := collection.Find(bson.D{{"attribute_id", Id}}).Count()
		lib.LogFatal(err)
		if n <= history {
			return nil
		}
		todel := KBHistory{}
		collection.Find(bson.D{{"attribute_id", Id}}).Sort("when").One(&todel)
		if todel.Id != "" {
			collection.RemoveId(todel.Id)
		} else {
			return nil
		}
	}
	return nil
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
