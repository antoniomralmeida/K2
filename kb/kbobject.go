package kb

import (
	"encoding/json"
	"log"

	"github.com/antoniomralmeida/k2/initializers"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (o *KBObject) Persist() error {
	collection := initializers.GetDb().C("KBObject")
	if o.Id == "" {
		o.Id = bson.NewObjectId()
		return collection.Insert(o)
	} else {
		return collection.UpdateId(o.Id, o)
	}
}

func (o *KBObject) String() string {
	j, err := json.MarshalIndent(*o, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}

func (o *KBObject) Delete() error {
	//TODO: Reiniciar KB
	return nil
}

func FindAllObjects(filter bson.M, sort string, os *[]KBObject) error {
	collection := initializers.GetDb().C("KBObject")
	idx, _ := collection.Indexes()
	if len(idx) == 1 {
		err := collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(filter).Sort(sort).All(os)
}

func (o *KBObject) GetWorkspaces() (ret []*KBWorkspace) {
	for i := range GKB.Workspaces {
		for j := range GKB.Workspaces[i].Objects {
			if GKB.Workspaces[i].Objects[j].KBObject == o {
				ret = append(ret, &GKB.Workspaces[i])
			}
		}
	}
	return
}
