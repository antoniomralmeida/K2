package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (w *KBWorkspace) Persist() error {
	collection := db.GetDb().C("KBWorkspace")
	if w.Id == "" {
		w.Id = bson.NewObjectId()
		return collection.Insert(w)
	} else {
		return collection.UpdateId(w.Id, w)
	}
}

func FindAllWorkspaces(sort string, ws *[]KBWorkspace) error {
	collection := db.GetDb().C("KBWorkspace")
	idx, err := collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"workspace"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(bson.M{}).Sort(sort).All(ws)
}
