package kb

import (
	"log"
	"time"

	"github.com/antoniomralmeida/k2/db"
	"gopkg.in/mgo.v2/bson"
)

func (r *KBRule) Run() {

	log.Println("run...", r.Rule)
	for _, x := range r.bin {
		switch x.typebin {
		case b_initially:
			{

			}
		}
	}
	r.lastexecution = time.Now()
}

func (r *KBRule) Persist() error {
	collection := db.GetDb().C("KBRule")
	if r.Id == "" {
		r.Id = bson.NewObjectId()
		return collection.Insert(r)
	} else {
		return collection.UpdateId(r.Id, r)
	}
}

func FindAllRules(sort string, rs *[]KBRule) error {
	collection := db.GetDb().C("KBRule")
	return collection.Find(bson.M{}).Sort(sort).All(rs)
}
