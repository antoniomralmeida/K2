package initializers

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/antoniomralmeida/k2/lib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
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
	Log(err, Fatal)
	//ping db
	err = client.Ping(ctx, readpref.Primary())
	Log(err, Fatal)
	db = client.Database(dbName)
}

func GetCollection(name string) (context.Context, *mongo.Collection) {
	return ctx, db.Collection(name)
}

// CreateUniqueIndex create UniqueIndex
func CreateUniqueIndex(collection string, keys ...string) {
	keysDoc := bsonx.Doc{}
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	idxRet, err := db.Collection(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		options.CreateIndexes().SetMaxTime(10*time.Second),
	)
	if err != nil {
		Log(err, Fatal)
	}
	Log("collection.Indexes().CreateOne:"+idxRet, Info)
}
