package db

import (
	"log"

	"gopkg.in/mgo.v2"
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
