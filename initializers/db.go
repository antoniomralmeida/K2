package initializers

import (
	"context"
	"os"

	"github.com/antoniomralmeida/k2/lib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var ctx context.Context

func ConnectDB() {
	Log("ConnectDB", Info)

	dsn := os.Getenv("DSN")
	dbName := os.Getenv("DB")
	Log(lib.Ping(dsn), Fatal)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		Log(err, Fatal)
	}
	db = client.Database(dbName)
}

func GetCollection(name string) (context.Context, *mongo.Collection) {
	return ctx, db.Collection(name)
}
