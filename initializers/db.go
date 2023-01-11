package initializers

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/antoniomralmeida/k2/lib"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func ConnectDB() {
	Log("ConnectDB", Info)

	dsn := os.Getenv("DSN")
	dbName := os.Getenv("DB")
	Log(lib.Ping(dsn), Fatal)

	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(dsn))
	Log(err, Fatal)
	client, err := mgm.NewClient()
	Log(err, Fatal)
	defer client.Disconnect(mgm.Ctx())
	//ping db
	err = client.Ping(mgm.Ctx(), readpref.Primary())
	Log(err, Fatal)
}

// CreateUniqueIndex create UniqueIndex
func CreateUniqueIndex(coll *mgm.Collection, keys ...string) {
	keysDoc := bsonx.Doc{}
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	idxRet, err := coll.Indexes().CreateOne(
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
