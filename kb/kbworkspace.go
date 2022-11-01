package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/initializers"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (w *KBWorkspace) Persist() error {
	collection := initializers.GetDb().C("KBWorkspace")
	if w.Id == "" {
		w.Id = bson.NewObjectId()
		return collection.Insert(w)
	} else {
		return collection.UpdateId(w.Id, w)
	}
}

func FindAllWorkspaces(sort string, ws *[]KBWorkspace) error {
	collection := initializers.GetDb().C("KBWorkspace")
	idx, _ := collection.Indexes()
	if len(idx) == 1 {
		err := collection.EnsureIndex(mgo.Index{Key: []string{"workspace"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(bson.M{}).Sort(sort).All(ws)
}
