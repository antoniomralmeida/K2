package kb

import (
	"github.com/antoniomralmeida/k2/db"
	"gopkg.in/mgo.v2/bson"
)

func (m *KBMutex) Persist() error {
	collection := db.GetDb().C("KBMutex")
	if m.Id == "" {
		m.Id = bson.NewObjectId()
		return collection.Insert(m)
	} else {
		return collection.UpdateId(m.Id, m)
	}
}

func (m *KBMutex) FindOne() error {
	collection := db.GetDb().C("KBMutex")
	return collection.Find(bson.D{}).One(m)
}
