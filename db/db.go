package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func GetDb() *mgo.Database {
	return db
}

func ConnectDB(uri string, dbName string) {
	log.Println("ConnectDB")

	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(dbName)
}
