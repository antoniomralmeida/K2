package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (c *KBClass) addAntecedentRules(r *KBRule) {
	found := false
	for i, _ := range c.antecedentRules {
		if c.antecedentRules[i] == r {
			found = true
			break
		}
	}
	if !found {
		c.antecedentRules = append(c.antecedentRules, r)
	}
}

func (c *KBClass) addConsequentRules(r *KBRule) {
	found := false
	for i, _ := range c.consequentRules {
		if c.consequentRules[i] == r {
			found = true
			break
		}
	}
	if !found {
		c.consequentRules = append(c.consequentRules, r)
	}
}

func FindAllClasses(sort string, cs *[]KBClass) error {
	collection := db.GetDb().C("KBClass")
	idx, err := collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	return collection.Find(bson.M{}).Sort(sort).All(cs)
}

func (class *KBClass) Persist() error {
	collection := db.GetDb().C("KBClass")
	if class.Id == "" {
		class.Id = bson.NewObjectId()
		return collection.Insert(class)
	} else {
		return collection.UpdateId(class.Id, class)
	}
}

func (class *KBClass) FindOne(p bson.D) error {
	collection := db.GetDb().C("KBClass")
	return collection.Find(p).One(class)
}
