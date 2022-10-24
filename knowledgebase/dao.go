package knowledgebase

import (
	"github.com/antoniomralmeida/k2/knowledgebase"
	"gopkg.in/mgo.v2/bson"
)

func Persist(class *knowledgebase.KBClass) error {
	collection := db.db_write.C("KBClass")
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}

func FindAll(classes *[]knowledgebase.KBClass, sort string) error {
	collection := db.db_write.C("KBClass")
	return collection.Find(bson.M{}).Sort(sort).All(classes)
}
