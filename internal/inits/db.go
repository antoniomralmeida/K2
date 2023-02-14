package inits

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/pkg/dsn"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func ConnectDB() {

	type Ping struct {
		mgm.DefaultModel `json:",inline" bson:",inline"`
	}
	Log("ConnectDB", Info)

	server := os.Getenv("MONGO_SERVER")
	dsn := dsn.Decode(server)
	Log(dsn, Info)
	Log(lib.Ping(server), Fatal)
	err := mgm.SetDefaultConfig(nil, dsn.Query("database"), options.Client().ApplyURI(server))
	Log(err, Fatal)
	client, err := mgm.NewClient()
	Log(err, Fatal)
	defer client.Disconnect(mgm.Ctx())
	//ping db
	ping := new(Ping)
	err = mgm.Coll(ping).FindByID(0, ping)
	if err != mongo.ErrNoDocuments {
		Log(err, Fatal)
	}
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

type ModelExt interface {
	// PrepareID converts the id value if needed, then
	// returns it (e.g convert string to objectId).
	PrepareID(id interface{}) (interface{}, error)

	GetID() interface{}
	SetID(id interface{})
	GetPrimitiveUpdateAt() primitive.DateTime
}

func Persist(model ModelExt) error {
	if !primitive.IsValidObjectID(fmt.Sprint(model.GetID())) {
		return mgm.Coll(model).Create(model)
	} else {
		res := mgm.Coll(model).FindOne(mgm.Ctx(), bson.D{{Key: "_id", Value: model.GetID()}, {Key: "updated_at", Value: model.GetPrimitiveUpdateAt()}})
		if res.Err() == nil {
			return mgm.Coll(model).Update(model)
		} else {
			return res.Err()
		}
	}
}
