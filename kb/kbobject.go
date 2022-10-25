package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (class *KBObject) Persist() error {
	collection := db.GetDb().C("KBObject")
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}

func FindAllObjects(sort string, os *[]KBObject) error {
	collection := db.GetDb().C("KBObject")
	idx, err := collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(bson.M{}).Sort(sort).All(os)
}
