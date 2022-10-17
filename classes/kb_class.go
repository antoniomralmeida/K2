package classes

import (
	"log"
	"sort"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type KBAttributeType string

const (
	KBString KBAttributeType = "String"
	KBDate                   = "Date"
	KBNumber                 = "Number"
	KBList                   = "List"
)

type KBSource string

const (
	User       KBSource = "User"
	PLC                 = "PLC"
	History             = "History"
	Simulation          = "Simulation"
)

type KBSimulation string

const (
	Default       KBSimulation = ""
	MonteCarlo                 = "Monte Carlo"
	MovingAverage              = "Moving Average"
	Interpolation              = "interpolation"
)

type KBAttribute struct {
	Id               bson.ObjectId   `bson:"id,omitempty"`
	Name             string          `bson:"name"`
	AType            KBAttributeType `bson:"atype"`
	Options          []string        `bson:"options,omitempty"`
	Sources          []KBSource      `bson:"sources"`
	KeepHistory      int             `bson:"keephistory"`
	ValidityInterval int             `bson:"validityinterval"`
	Deadline         int             `bson:"deadline"`
	Simulation       KBSimulation    `bson:"simulation,omitempty"`
}

type KBClass struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Icon        string        `bson:"icon"`
	Parent      bson.ObjectId `bson:"parent_id,omitempty"`
	ParentClass *KBClass      `bson:"-"`
	Attributes  []KBAttribute `bson:"attributes"`
}

type KBRule struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	Rule string        `bson:"rule"`
	bin  []*Token      `bson:"-"`
}

type KnowledgeBase struct {
	Classes    []KBClass
	Rules      []KBRule
	Workspaces []KBWorkspace
	ebnf       *EBNF
	db         *mgo.Database
}

type KBHistory struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Attribute bson.ObjectId `bson:"object_id"`
	When      time.Time     `bson:"when"`
	Value     string        `bson:"value"`
	Certainty float64       `bson:"certainty"`
	Source    KBSource      `bson:"source"`
}

type KBAttributeObject struct {
	Id          bson.ObjectId `bson:"id"`
	Attribute   bson.ObjectId `bson:"attribute_id"`
	KbAttribute *KBAttribute  `bson:"-"`
	KbHistory   *KBHistory    `bson:"-"`
}

type KBObject struct {
	Name       string              `bson:"name"`
	Class      bson.ObjectId       `bson:"class_id"`
	Top        int                 `bson:"top"`
	Left       int                 `bson:"left"`
	Bkclass    *KBClass            `bson:"-"`
	Attributes []KBAttributeObject `bson:"attributes"`
}

type KBWorkspace struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Workspace       string        `bson:"workspace"`
	Top             int           `bson:"top"`
	Left            int           `bson:"left"`
	Width           int           `bson:"width"`
	Height          int           `bson:"height"`
	BackgroundImage string        `bson:"backgroundimage,omitempty"`
	Objects         []KBObject    `bson:"objects"`
}

func (kb *KnowledgeBase) ConnectDB(uri string, dbName string) {
	log.Println("ConnectDB")
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}
	kb.db = session.DB(dbName)
}

/*
	func (kb *KnowledgeBase) sortClass() {
		//Sort Classes by Id
		sort.Slice(kb.Classes, func(i, j int) bool {
			return kb.Classes[i].Id > kb.Classes[j].Id
		})
	}
*/
func (kb *KnowledgeBase) findClassBin(id bson.ObjectId, i int, j int) int {
	if j >= i {
		avg := (i + j) / 2
		if kb.Classes[avg].Id == id {
			return avg
		} else if kb.Classes[avg].Id < id {
			return kb.findClassBin(id, i, avg-1)
		} else {
			return kb.findClassBin(id, avg+1, j)
		}
	} else {
		return -1
	}
}

func (kb *KnowledgeBase) NewClass(c *KBClass) {
	c.Id = bson.NewObjectId()
	for i, _ := range c.Attributes {
		c.Attributes[i].Id = bson.NewObjectId()
	}
	collection := kb.db.C("Class")
	err := collection.Insert(c)
	if err != nil {
		log.Fatal(err)
	}
	kb.Classes = append(kb.Classes, *c)
	err = collection.Find(bson.M{}).Sort("_id").All(&kb.Classes)
	if err != nil {
		log.Fatal(err)
	}
}

func (kb *KnowledgeBase) UpdateClass(c *KBClass) {
	collection := kb.db.C("Class")
	for i, _ := range c.Attributes {
		if c.Attributes[i].Id == "" {
			c.Attributes[i].Id = bson.NewObjectId()
		}
	}
	err := collection.UpdateId(c.Id, c)
	if err != nil {
		log.Fatal(err)
	}
}

func (kb *KnowledgeBase) NewWorkspace(w *KBWorkspace) {
	w.Id = bson.NewObjectId()
	collection := kb.db.C("Workspace")
	err := collection.Insert(w)
	if err != nil {
		log.Fatal(err)
	}
	kb.Workspaces = append(kb.Workspaces, *w)
}

func (kb *KnowledgeBase) UpdateWorkspace(w *KBWorkspace) {
	collection := kb.db.C("Workspace")
	err := collection.UpdateId(w.Id, w)
	if err != nil {
		log.Fatal(err)
	}
}

func (kb *KnowledgeBase) NewObject(w int, o *KBObject, c *KBClass) {
	o.Class = c.Id
	kb.Workspaces[w].Objects = append(kb.Workspaces[w].Objects, *o)
	oi := len(kb.Workspaces[w].Objects) - 1
	for _, x := range kb.FindAttributes(c) {
		n := KBAttributeObject{Id: bson.NewObjectId(), Attribute: x.Id, KbAttribute: x}
		kb.Workspaces[w].Objects[oi].Attributes = append(kb.Workspaces[w].Objects[oi].Attributes, n)
	}
	kb.UpdateWorkspace(&kb.Workspaces[w])

}

func (kb *KnowledgeBase) SaveValue(w int, a *KBAttributeObject, h *KBHistory) {
	h.Id = bson.NewObjectId()
	h.Attribute = a.Id
	h.When = time.Now()
	collection := kb.db.C("History")
	err := collection.Insert(h)
	if err != nil {
		log.Fatal(err)
	}
	a.KbHistory = h
}

func (kb *KnowledgeBase) GetClass(name string) *KBClass {
	var ret KBClass
	collection := kb.db.C("Class")
	err := collection.Find(bson.D{{"name", name}}).One(&ret)
	if err != nil {
		log.Fatal(err)
	}
	return &ret
}

func (kb *KnowledgeBase) FindClass(id bson.ObjectId) int {
	return kb.findClassBin(id, 0, len(kb.Classes)-1)
}

func (kb *KnowledgeBase) FindWorkspaceByName(name string) int {
	for i, x := range kb.Workspaces {
		if x.Workspace == name {
			return i
		}
	}
	return -1
}

func (kb *KnowledgeBase) FindAttributes(c *KBClass) []*KBAttribute {
	var ret []*KBAttribute
	if c.ParentClass != nil {
		for _, x := range kb.FindAttributes(c.ParentClass) {
			ret = append(ret, x)
		}
	} else {
		for i, _ := range c.Attributes {
			ret = append(ret, &c.Attributes[i])
		}
	}
	return ret
}

func (kb *KnowledgeBase) NewRule(r *KBRule) {
	r.Id = bson.NewObjectId()
	collection := kb.db.C("Rule")
	err := collection.Insert(r)
	if err != nil {
		log.Fatal(err)
	}
}

func (kb *KnowledgeBase) ReadBK() {
	log.Println("ReadBK")
	collection := kb.db.C("Class")
	err := collection.Find(bson.M{}).Sort("_id").All(&kb.Classes)
	if err != nil {
		log.Fatal(err)
	}
	idx, err := collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name", "attributes.name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}

	for j, c := range kb.Classes {
		if c.Parent != "" {
			pc := kb.FindClass(c.Parent)
			if pc != -1 {
				kb.Classes[j].ParentClass = &kb.Classes[pc]
			} else {
				log.Fatal("Parent of Class " + c.Name + " not found!")
			}
		}
	}

	collection = kb.db.C("Rule")
	err = collection.Find(bson.M{}).Sort("_id").All(&kb.Rules)
	if err != nil {
		log.Fatal(err)
	}

	collection = kb.db.C("Workspace")
	history := kb.db.C("History")
	err = collection.Find(bson.M{}).All(&kb.Workspaces)
	if err != nil {
		log.Fatal(err)
	}
	idx, err = collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name", "objects.name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, w := range kb.Workspaces {
		for j, o := range w.Objects {
			c := kb.FindClass(o.Class)
			if c != -1 {
				kb.Workspaces[i].Objects[j].Bkclass = &kb.Classes[c]
				attrs := kb.FindAttributes(&kb.Classes[c])
				sort.Slice(attrs, func(i, j int) bool {
					return attrs[i].Id < attrs[j].Id
				})
				for k, x := range o.Attributes {
					for _, y := range attrs {
						if y.Id == x.Attribute {
							kb.Workspaces[i].Objects[j].Attributes[k].KbAttribute = y
							break
						}
						if y.Id > x.Attribute {
							break
						}
					}
					if kb.Workspaces[i].Objects[j].Attributes[i].KbAttribute == nil {
						log.Fatal("Attribute not found ", x.Attribute)
					}
					//Obter ultimo valor
					err = history.Find(bson.D{{"Attribute", x.Id}}).Sort("When").One(kb.Workspaces[i].Objects[j].Attributes[i].KbHistory)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				log.Fatal("Class of object " + o.Name + " not found!")
			}
		}
	}
}
