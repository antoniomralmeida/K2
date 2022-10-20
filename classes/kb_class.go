package classes

import (
	"log"
	"main/lib"
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
	Objects    []KBObject
	ebnf       *EBNF
	db         *mgo.Database
	stack      []*KBRule
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
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Name            string        `bson:"name"`
	Icon            string        `bson:"icon"`
	Parent          bson.ObjectId `bson:"parent_id,omitempty"`
	Attributes      []KBAttribute `bson:"attributes"`
	ParentClass     *KBClass      `bson:"-"`
	antecedentRules []*KBRule     `bson:"-"`
	consequentRules []*KBRule     `bson:"-"`
}

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

type KBRule struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Rule     string        `bson:"rule"`
	Priority byte          `bson:"priority"` //0..100
	bin      []*BIN        `bson:"-"`
}

func (kb *KnowledgeBase) linkerRule(r *KBRule, bin []*BIN) {
	// Find references of objects in KB
	log.Println("Linking Prodution Rule: ", r.Id)
	consequent := -1
	for j, x := range bin {
		switch x.typebin {
		case b_initially:
			kb.stack = append(kb.stack, r)
		case b_then:
			consequent = j
		}
		switch x.tokentype {
		case Object:
			if bin[j].object == nil {
				bin[j].object = kb.FindObjectByName(x.token)
			}
		case Class:
			if bin[j].class == nil {
				bin[j].class = kb.FindClassByName(x.token, true)
			}
		case Attribute:
			ref := -1
			if bin[j+1].typebin == b_of {
				ref = j + 3
			} else {
				for z := j - 1; z >= 0; z-- {
					if bin[z].tokentype == Object || bin[z].tokentype == Class {
						ref = z
						break
					}
				}
			}
			if ref != -1 {
				if bin[ref].tokentype == Object {
					if bin[ref].object == nil {
						bin[ref].object = kb.FindObjectByName(bin[ref].token)
					}
					bin[j].attribute = kb.FindAttributeObjectByClass(bin[ref].object, x.token)
					break
				} else if bin[ref].tokentype == Class {
					if bin[ref].class == nil {
						bin[j].class = kb.FindClassByName(x.token, true)
					}
					bin[j].attribute = kb.FindAttribute(bin[ref].class, x.token)
					break
				}
			} else {
				log.Fatal("Attribute not found in KB! ", x.token)
			}
		case DynamicReference:
			{
				if consequent == -1 {
					for z := j - 1; z >= 0; z-- {
						if bin[z].tokentype == Object {
							bin[j].object = bin[z].object
							break
						} else if bin[z].tokentype == Class {
							bin[j].class = bin[z].class
							break
						}
					}
				} else {
					for z := consequent - 1; z >= 0; z-- {
						if bin[z].tokentype == DynamicReference && bin[z].token == x.token {
							bin[j].object = bin[z].object
							bin[j].class = bin[z].class
							break
						}
					}
				}
			}

		case Constant:
			{
				if !lib.IsNumber(x.token) {
					ok := false
					for z := j - 1; z >= 0; z-- {
						if bin[z].tokentype == Attribute {
							if bin[z].attribute != nil {
								for _, o := range bin[z].attribute.Options {
									if x.token == o {
										ok = true
										break
									}
								}
							}
						}
					}
					if !ok {
						log.Fatal("Constant not found in KB! ", x.token)
					}
				}
			}
		}
		cl := bin[j].class
		if bin[j].object != nil {
			cl = bin[j].object.Bkclass
		}
		if cl != nil {
			if consequent != -1 {
				cl.addConsequentRules(r)
			} else {
				cl.addAntecedentRules(r)
			}
		}
		//fmt.Println(bin[j])
	}
	r.bin = bin
}
func (r *KBRule) Run() {
	for _, x := range r.bin {
		switch x.typebin {
		case b_initially:
			{

			}
		}
	}
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
	Id         bson.ObjectId       `bson:"_id"`
	Name       string              `bson:"name"`
	Class      bson.ObjectId       `bson:"class_id"`
	Top        int                 `bson:"top"`
	Left       int                 `bson:"left"`
	Attributes []KBAttributeObject `bson:"attributes"`
	Bkclass    *KBClass            `bson:"-"`
}

type KBWorkspace struct {
	Id              bson.ObjectId   `bson:"_id,omitempty"`
	Workspace       string          `bson:"workspace"`
	Top             int             `bson:"top"`
	Left            int             `bson:"left"`
	Width           int             `bson:"width"`
	Height          int             `bson:"height"`
	BackgroundImage string          `bson:"backgroundimage,omitempty"`
	Objects         []bson.ObjectId `bson:"objects"`
	KBObjects       []*KBObject     `bson:"-"`
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

func (kb *KnowledgeBase) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i, _ := range attrs {
		attrs[i].Id = bson.NewObjectId()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	collection := kb.db.C("Class")
	err := collection.UpdateId(c.Id, c)
	if err != nil {
		log.Fatal(err)
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

func (kb *KnowledgeBase) NewObject(c *KBClass, name string) *KBObject {
	collection := kb.db.C("Object")
	o := KBObject{Id: bson.NewObjectId(), Name: name, Class: c.Id, Bkclass: c}
	for _, x := range kb.FindAttributes(c) {
		n := KBAttributeObject{Id: bson.NewObjectId(), Attribute: x.Id, KbAttribute: x}
		o.Attributes = append(o.Attributes, n)
	}
	collection.Insert(&o)
	return &o
}

func (kb *KnowledgeBase) LinkObjects(ws *KBWorkspace, objs ...*KBObject) {
	for i, _ := range objs {
		ws.Objects = append(ws.Objects, objs[i].Id)
		ws.KBObjects = append(ws.KBObjects, objs[i])
	}
	kb.UpdateWorkspace(ws)
}

func (kb *KnowledgeBase) FindObjectByName(name string) *KBObject {
	return kb.findObjectByNameBin(name, 0, len(kb.Objects)-1)
}

func (kb *KnowledgeBase) findObjectByNameBin(name string, i int, j int) *KBObject {
	if j >= i {
		avg := (i + j) / 2
		if kb.Objects[avg].Name == name {
			return &kb.Objects[avg]
		} else if kb.Objects[avg].Name > name {
			return kb.findObjectByNameBin(name, i, avg-1)
		} else {
			return kb.findObjectByNameBin(name, avg+1, j)
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

func (kb *KnowledgeBase) FindClassByName(name string, mandatory bool) *KBClass {
	var ret KBClass
	collection := kb.db.C("Class")
	err := collection.Find(bson.D{{"name", name}}).One(&ret)
	if err != nil && mandatory {
		log.Fatal(err)
	}
	return kb.FindClassById(ret.Id)
}

func (kb *KnowledgeBase) FindClassById(id bson.ObjectId) *KBClass {
	return kb.findClassBin(id, 0, len(kb.Classes)-1)
}

func (kb *KnowledgeBase) FindAttribute(c *KBClass, name string) *KBAttribute {
	attrs := kb.FindAttributes(c)
	for i, x := range attrs {
		if x.Name == name {
			return attrs[i]
		}
	}
	return nil
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

func (kb *KnowledgeBase) NewAttributeObject(obj *KBObject, attr *KBAttribute) *KBAttributeObject {
	a := KBAttributeObject{Attribute: attr.Id, Id: bson.NewObjectId()}
	obj.Attributes = append(obj.Attributes, a)
	collection := kb.db.C("Object")
	err := collection.UpdateId(obj.Id, obj)
	if err != nil {
		log.Fatal(err)
	}
	return &a
}

func (kb *KnowledgeBase) FindAttributeObjectByClass(obj *KBObject, attr string) *KBAttribute {
	for i, _ := range obj.Bkclass.Attributes {
		if obj.Bkclass.Attributes[i].Name == attr {
			return &obj.Bkclass.Attributes[i]
		}
	}
	return nil
}

func (kb *KnowledgeBase) NewRule(priority byte, rule string) *KBRule {
	_, bin, err := kb.ebnf.Parsing(kb, rule)
	if err != nil {
		log.Fatal(err)
	}
	r := KBRule{Rule: rule, Priority: priority}
	r.Id = bson.NewObjectId()
	kb.linkerRule(&r, bin)
	collection := kb.db.C("Rule")
	err = collection.Insert(r)
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

	collection = kb.db.C("Object")
	err = collection.Find(bson.M{}).Sort("name").All(&kb.Objects)
	if err != nil {
		log.Fatal(err)
	}
	idx, err = collection.Indexes()
	if len(idx) == 1 {
		err = collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
		if err != nil {
			log.Fatal(err)
		}
	}
	history := kb.db.C("History")
	for j, o := range kb.Objects {
		c := kb.FindClassById(o.Class)
		if c != nil {
			kb.Objects[j].Bkclass = c
			attrs := kb.FindAttributes(c)
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].Id < attrs[j].Id
			})
			for k, x := range o.Attributes {
				for l, y := range attrs {
					if y.Id == x.Attribute {
						kb.Objects[j].Attributes[k].KbAttribute = attrs[l]
						break
					}
					if y.Id > x.Attribute {
						break
					}
				}
				if kb.Objects[j].Attributes[k].KbAttribute == nil {
					log.Fatal("Attribute not found ", x.Attribute)
				}
				//Obter ultimo valor
				h := KBHistory{}
				err = history.Find(bson.D{{"attribute_id", x.Id}}).Sort("-when").One(&h)
				if err != nil {
					if err.Error() != "not found" {
						log.Println(err)
					}
				} else {
					kb.Objects[j].Attributes[k].KbHistory = &h
				}
			}
		} else {
			log.Fatal("Class of object " + o.Name + " not found!")
		}
	}

	collection = kb.db.C("Workspace")
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
	}

	collection = kb.db.C("Rule")
	err = collection.Find(bson.M{}).Sort("-priority").All(&kb.Rules)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range kb.Rules {
		_, bin, err := kb.ebnf.Parsing(kb, kb.Rules[i].Rule)
		if err != nil {
			log.Fatal(err)
		}
		kb.linkerRule(&kb.Rules[i], bin)
	}
}

func (kb *KnowledgeBase) ReadEBNF(file string) {
	ebnf := EBNF{}
	kb.ebnf = &ebnf
	kb.ebnf.ReadToken(file)
}

func (kb *KnowledgeBase) PrintEBNF() {
	kb.ebnf.PrintEBNF()
}
