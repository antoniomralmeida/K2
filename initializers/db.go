package initializers

import (
	"os"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func GetDb() *mgo.Database {
	return db
}

func ConnectDB() {
	Log("ConnectDB", Info)

	dsn := os.Getenv("DSN")
	dbName := os.Getenv("db")

	session, err := mgo.Dial(dsn)
	if err != nil {
		Log(err, Fatal)
	}
	db = session.DB(dbName)
}
