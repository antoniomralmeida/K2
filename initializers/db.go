package initializers

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var ctx context.Context

/*
func GetDb() *mongo.Database {
	return db
}
*/

func ConnectDB() {
	Log("ConnectDB", Info)

	dsn := os.Getenv("DSN")
	dbName := os.Getenv("db")
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		Log(err, Fatal)
	}
	db = client.Database(dbName)
}

func GetCollection(name string) (context.Context, *mongo.Collection) {
	return ctx, db.Collection(name)
}

/*
func Insert(c *mongo.Collection, docs ...interface{}) {
	var fiels []string
	var values []any
	for d := range docs {
		v := reflect.ValueOf(d).Elem()
		for j := 0; j < v.NumField(); j++ {
			fiels = append(fiels, v.Type().Field(j).Name)
			values = append(values, v.Field(j))
		}
		c.InsertOne(ctx, bson)
	}
}
*/
