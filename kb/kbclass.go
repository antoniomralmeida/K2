package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindAllClasses(sort string, cs *[]KBClass) error {
	collection := db.GetDb().C("KBClass")
	idx, _ := collection.Indexes()
	if len(idx) == 1 {
		err := collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(bson.M{}).Sort(sort).All(cs)
}

func (class *KBClass) Persist() error {
	collection := db.GetDb().C("KBClass")
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}

func (class *KBClass) FindOne(p bson.D) error {
	collection := db.GetDb().C("KBClass")
	return collection.Find(p).One(class)
}

func (class *KBClass) Delete() error {
	//TODO: Restart KB
	return nil
}
