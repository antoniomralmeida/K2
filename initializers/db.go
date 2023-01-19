package initializers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func ConnectDB() {
	Log("ConnectDB", Info)

	dsn := os.Getenv("DSN")
	dbName := os.Getenv("DB")

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

/*
	func (u *User) GetPrimitiveUpdateAt() primitive.DateTime {
		return primitive.NewDateTimeFromTime(u.UpdatedAt)
	}
*/
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
		res := mgm.Coll(model).FindOne(mgm.Ctx(), bson.D{{"_id", model.GetID()}, {"updated_at", model.GetPrimitiveUpdateAt()}})
		if res.Err() == nil {
			return mgm.Coll(model).Update(model)
		} else {
			return res.Err()
		}
	}
}
