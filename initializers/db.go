package initializers

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func GetDb() *mgo.Database {
	return db
}

func ConnectDB() {
	log.Println("ConnectDB")

	dsn := os.Getenv("DSN")
	dbName := os.Getenv("db")

	session, err := mgo.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(dbName)
}
