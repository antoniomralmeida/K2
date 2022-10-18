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

type KnowledgeBase struct {
	Classes    []KBClass
	Rules      []KBRule
	Workspaces []KBWorkspace
	ebnf       *EBNF
	db         *mgo.Database
}

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

type KBHistory struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Attribute bson.ObjectId `bson:"attribute_id"`
	When      time.Time     `bson:"when"`
	Value     string        `bson:"value"`
	Certainty float64       `bson:"certainty,omitempty"`
	Source    KBSource      `bson:"source"`
}

type KBAttributeObject struct {
	Id          bson.ObjectId `bson:"id"`
	Attribute   bson.ObjectId `bson:"attribute_id"`
	KbAttribute *KBAttribute  `bson:"-"`
	KbHistory   *KBHistory    `bson:"-"`
}

func (ao *KBAttributeObject) Value() string {
	if ao.KbHistory != nil {
		return ao.KbHistory.Value
	} else {
		return ""
	}
}

type KBObject struct {
	Id         bson.ObjectId       `bson:"id"`
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

func (kb *KnowledgeBase) findClassBin(id bson.ObjectId, i int, j int) *KBClass {
	if j >= i {
		avg := (i + j) / 2
		if kb.Classes[avg].Id == id {
			return &kb.Classes[avg]
		} else if kb.Classes[avg].Id < id {
			return kb.findClassBin(id, i, avg-1)
		} else {
			return kb.findClassBin(id, avg+1, j)
		}
	} else {
		return nil
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

func (kb *KnowledgeBase) NewWorkspace(name string, icone string) *KBWorkspace {
	w := KBWorkspace{Workspace: name, BackgroundImage: icone}
	w.Id = bson.NewObjectId()
	collection := kb.db.C("Workspace")
	err := collection.Insert(w)
	if err != nil {
		log.Fatal(err)
	}
	kb.Workspaces = append(kb.Workspaces, w)
	return &w
}

func (kb *KnowledgeBase) UpdateWorkspace(w *KBWorkspace) {
	collection := kb.db.C("Workspace")
	err := collection.UpdateId(w.Id, w)
	if err != nil {
		log.Fatal(err)
	}
}

func (kb *KnowledgeBase) FindWorkspaceByName(name string) *KBWorkspace {
	for i, _ := range kb.Workspaces {
		if kb.Workspaces[i].Workspace == name {
			return &kb.Workspaces[i]
		}
	}
	log.Fatal("Workspace not found!")
	return nil
}

func (kb *KnowledgeBase) NewObject(ws *KBWorkspace, c *KBClass, name string) *KBObject {
	o := KBObject{Id: bson.NewObjectId(), Name: name, Class: c.Id, Bkclass: c}
	for _, x := range kb.FindAttributes(c) {
		n := KBAttributeObject{Id: bson.NewObjectId(), Attribute: x.Id, KbAttribute: x}
		o.Attributes = append(o.Attributes, n)
	}
	ws.Objects = append(ws.Objects, o)
	kb.UpdateWorkspace(ws)
	return &o
}

func (kb *KnowledgeBase) FindObjectByName(ws *KBWorkspace, name string) *KBObject {
	return kb.findObjectByNameBin(ws, name, 0, len(ws.Objects)-1)
}

func (kb *KnowledgeBase) findObjectByNameBin(ws *KBWorkspace, name string, i int, j int) *KBObject {
	if j >= i {
		avg := (i + j) / 2
		if ws.Objects[avg].Name == name {
			return &ws.Objects[avg]
		} else if ws.Objects[avg].Name > name {
			return kb.findObjectByNameBin(ws, name, i, avg-1)
		} else {
			return kb.findObjectByNameBin(ws, name, avg+1, j)
		}
	} else {
		return nil
	}
}
func (kb *KnowledgeBase) SaveValue(attr *KBAttributeObject, value string, source KBSource) *KBHistory {
	if attr != nil {
		h := KBHistory{Id: bson.NewObjectId(), Attribute: attr.Id, When: time.Now(), Value: value, Source: source}

		collection := kb.db.C("History")
		err := collection.Insert(h)
		if err != nil {
			log.Fatal(err)
		}
		attr.KbHistory = &h
		return &h
	} else {
		log.Fatal("Invalid Attribute of Object!")
		return nil
	}
}

func (kb *KnowledgeBase) FindClassByName(name string) *KBClass {
	var ret KBClass
	collection := kb.db.C("Class")
	err := collection.Find(bson.D{{"name", name}}).One(&ret)
	if err != nil {
		log.Fatal(err)
	}
	return kb.FindClassById(ret.Id)
}

func (kb *KnowledgeBase) FindClassById(id bson.ObjectId) *KBClass {
	return kb.findClassBin(id, 0, len(kb.Classes)-1)
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

func (kb *KnowledgeBase) FindAttributeObject(obj *KBObject, attr string) *KBAttributeObject {
	for i, _ := range obj.Attributes {
		if obj.Attributes[i].KbAttribute.Name == attr {
			return &obj.Attributes[i]
		}
	}
	return nil
}

func (kb *KnowledgeBase) NewRule(rule string) *KBRule {
	r := KBRule{Rule: rule}
	r.Id = bson.NewObjectId()
	collection := kb.db.C("Rule")
	err := collection.Insert(r)
	if err != nil {
		log.Fatal(err)
	}
	kb.Rules = append(kb.Rules, r)
	return &r
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
		log.Println("Prepare Class ", c.Name)
		if c.Parent != "" {
			pc := kb.FindClassById(c.Parent)
			if pc != nil {
				kb.Classes[j].ParentClass = pc
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
		err = collection.EnsureIndex(mgo.Index{Key: []string{"workspace"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
		err = collection.EnsureIndex(mgo.Index{Key: []string{"workspace", "objects.name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, w := range kb.Workspaces {
		log.Println("Prepare Workspace ", w.Workspace)
		for j, o := range w.Objects {
			c := kb.FindClassById(o.Class)
			if c != nil {
				kb.Workspaces[i].Objects[j].Bkclass = c
				attrs := kb.FindAttributes(c)
				sort.Slice(attrs, func(i, j int) bool {
					return attrs[i].Id < attrs[j].Id
				})
				for k, x := range o.Attributes {
					for l, y := range attrs {
						if y.Id == x.Attribute {
							kb.Workspaces[i].Objects[j].Attributes[k].KbAttribute = attrs[l]
							break
						}
						if y.Id > x.Attribute {
							break
						}
					}
					if kb.Workspaces[i].Objects[j].Attributes[k].KbAttribute == nil {
						log.Fatal("Attribute not found ", x.Attribute)
					}
					//Obter ultimo valor
					h := KBHistory{}
					err = history.Find(bson.D{{"attribute_id", x.Id}}).Sort("-when").One(&h)
					if err != nil {
						log.Println(err)
					} else {
						kb.Workspaces[i].Objects[j].Attributes[k].KbHistory = &h
					}
				}
			} else {
				log.Fatal("Class of object " + o.Name + " not found!")
			}
		}
		sort.Slice(kb.Workspaces[i].Objects, func(a, b int) bool {
			return kb.Workspaces[i].Objects[a].Name < kb.Workspaces[i].Objects[b].Name
		})
	}
}

func (kb *KnowledgeBase) ReadEBNF(file string) {
	ebnf := EBNF{}
	kb.ebnf = &ebnf
	kb.ebnf.ReadToken(file)
}
