package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (o *KBObject) Persist() error {
	collection := db.GetDb().C("KBObject")
	if o.Id == "" {
		o.Id = bson.NewObjectId()
		return collection.Insert(o)
	} else {
		return collection.UpdateId(o.Id, o)
	}
}

func FindAllObjects(filter bson.M, sort string, os *[]KBObject) error {
	collection := db.GetDb().C("KBObject")
	idx, err := collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(filter).Sort(sort).All(os)
}
