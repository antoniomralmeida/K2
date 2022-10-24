package db

import (
	"log"
	"reflect"

	"github.com/antoniomralmeida/k2/knowledgebase"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	db_read  *mgo.Database
	db_write *mgo.Database
}

var db DB

func ConnectDB(uri string, dbName string) {
	log.Println("ConnectDB")

	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}
	db.db_read = session.DB(dbName)
	db.db_write = session.DB(dbName)

}

func Persist(class *knowledgebase.KBClass) error {
	collection := db.db_write.C("KBClass")
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}

func Find(class []*knowledgebase.KBClass) error {
	collection := db.db_write.C(reflect.TypeOf(*class).String())
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}
