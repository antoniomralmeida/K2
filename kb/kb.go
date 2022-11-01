package kb

import (
	"log"
	"sort"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"gopkg.in/mgo.v2/bson"
)

func (kb *KnowledgeBase) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i := range attrs {
		attrs[i].Id = bson.NewObjectId()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	lib.LogFatal(c.Persist())
}

func (kb *KnowledgeBase) NewClass(c *KBClass) {
	for i := range c.Attributes {
		c.Attributes[i].Id = bson.NewObjectId()
	}
	err := c.Persist()
	if err != nil {
		log.Fatal(err)
	}
	kb.Classes = append(kb.Classes, *c)
	kb.IdxClasses[c.Id] = c
}

func (kb *KnowledgeBase) UpdateClass(c *KBClass) {
	for i := range c.Attributes {
		if c.Attributes[i].Id == "" {
			c.Attributes[i].Id = bson.NewObjectId()
		}
	}
	lib.LogFatal(c.Persist())
}

func (kb *KnowledgeBase) NewWorkspace(name string, icone string) *KBWorkspace {
	w := KBWorkspace{Workspace: name, BackgroundImage: icone}
	log.Fatal(w.Persist())
	kb.Workspaces = append(kb.Workspaces, w)
	return &w
}

func (kb *KnowledgeBase) UpdateWorkspace(w *KBWorkspace) {
	lib.LogFatal(w.Persist())
}

func (kb *KnowledgeBase) FindWorkspaceByName(name string) *KBWorkspace {
	for i := range kb.Workspaces {
		if kb.Workspaces[i].Workspace == name {
			return &kb.Workspaces[i]
		}
	}
	log.Fatal("Workspace not found!")
	return nil
}

func (kb *KnowledgeBase) NewObject(c *KBClass, name string) *KBObject {

	o := KBObject{Name: name, Class: c.Id, Bkclass: c}
	for _, x := range kb.FindAttributes(c) {
		n := KBAttributeObject{Id: bson.NewObjectId(), Attribute: x.Id, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
	}
	lib.LogFatal(o.Persist())
	return &o
}

func (kb *KnowledgeBase) LinkObjects(ws *KBWorkspace, obj *KBObject, left int, top int) {
	ows := KBObjectWS{Object: obj.Id, Left: left, Top: top, KBObject: obj}
	ws.Objects = append(ws.Objects, ows)
	kb.UpdateWorkspace(ws)
}

func (kb *KnowledgeBase) FindObjectByName(name string) *KBObject {
	return kb.IdxObjects[name]
}

func (kb *KnowledgeBase) FindClassByName(nm string, mandatory bool) *KBClass {
	var ret KBClass
	err := ret.FindOne(bson.D{{"name", nm}})
	if err != nil && mandatory {
		log.Fatal(err)
	}
	return kb.IdxClasses[ret.Id]
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
	if c != nil {
		if c.ParentClass != nil {
			ret = append(ret, kb.FindAttributes(c.ParentClass)...)
		} else {
			for i := range c.Attributes {
				ret = append(ret, &c.Attributes[i])
			}
		}
	}
	return ret
}

func (kb *KnowledgeBase) FindAttributeObject(obj *KBObject, attr string) *KBAttributeObject {
	for i := range obj.Attributes {
		if obj.Attributes[i].KbAttribute.Name == attr {
			return &obj.Attributes[i]
		}
	}
	return nil
}

func (kb *KnowledgeBase) NewAttributeObject(obj *KBObject, attr *KBAttribute) *KBAttributeObject {
	a := KBAttributeObject{Attribute: attr.Id, Id: bson.NewObjectId()}
	obj.Attributes = append(obj.Attributes, a)
	log.Fatal(obj.Persist())
	return &a
}

func (kb *KnowledgeBase) NewRule(rule string, priority byte, interval int) *KBRule {
	_, bin, err := kb.ParsingCommand(rule)
	lib.LogFatal(err)
	r := KBRule{Rule: rule, Priority: priority, ExecutionInterval: interval}
	lib.LogFatal(r.Persist())
	kb.linkerRule(&r, bin)
	kb.Rules = append(kb.Rules, r)
	return &r
}
func (kb *KnowledgeBase) UpdateKB(name string, iotapi string) error {
	kb.Name = name
	kb.IOTApi = iotapi
	return kb.Persist()
}

func (kb *KnowledgeBase) Init() {
	log.Println("Init KB")

	kb.FindOne()
	if kb.Name == "" {
		kb.Name = "K2 System KB"
	}
	kb.Persist()
	kb.IdxClasses = make(map[bson.ObjectId]*KBClass)
	kb.IdxObjects = make(map[string]*KBObject)
	kb.IdxAttributeObjects = make(map[string]*KBAttributeObject)

	ebnf := ebnf.EBNF{}
	kb.ebnf = &ebnf
	kb.ebnf.ReadToken("./ebnf/k2.ebnf")

	FindAllClasses("_id", &kb.Classes)
	for j := range kb.Classes {
		kb.IdxClasses[kb.Classes[j].Id] = &kb.Classes[j]
	}

	for j, c := range kb.Classes {
		log.Println("Prepare Class ", c.Name)
		if c.Parent != "" {
			pc := kb.IdxClasses[c.Parent]
			if pc != nil {
				kb.Classes[j].ParentClass = pc
			} else {
				log.Fatal("Parent of Class " + c.Name + " not found!")
			}
		}
	}

	FindAllObjects(bson.M{}, "name", &kb.Objects)
	for j, o := range kb.Objects {
		kb.IdxObjects[o.Name] = &kb.Objects[j]
		c := kb.IdxClasses[o.Class]
		if c != nil {
			kb.Objects[j].Bkclass = c
			attrs := kb.FindAttributes(c)
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].Id < attrs[j].Id
			})
			for k, x := range o.Attributes {
				kb.Objects[j].Attributes[k].KbObject = &kb.Objects[j]
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
				kb.IdxAttributeObjects[o.Name+"."+kb.Objects[j].Attributes[k].KbAttribute.Name] = &kb.Objects[j].Attributes[k]

				//Obter ultimo valor
				h := KBHistory{}
				err := h.FindLast(bson.D{{"attribute_id", x.Id}})
				if err != nil {
					if err.Error() != "not found" {
						log.Println(err)
					}
				} else {
					kb.Objects[j].Attributes[k].KbHistory = &h
				}
				kb.Objects[j].Attributes[k].Validity()
			}
		} else {
			log.Fatal("Class of object " + o.Name + " not found!")
		}
	}

	FindAllWorkspaces("name", &kb.Workspaces)

	FindAllRules("_id", &kb.Rules)

	for i := range kb.Rules {
		_, bin, err := kb.ParsingCommand(kb.Rules[i].Rule)
		if err != nil {
			log.Fatal(err)
		}
		kb.linkerRule(&kb.Rules[i], bin)
	}
}

func (kb *KnowledgeBase) PrintEBNF() {
	kb.ebnf.PrintEBNF()
}

func (kb *KnowledgeBase) Persist() error {
	collection := initializers.GetDb().C("KnowledgeBase")
	if kb.Id == "" {
		kb.Id = bson.NewObjectId()
		return collection.Insert(kb)
	} else {
		return collection.UpdateId(kb.Id, kb)
	}
}

func (kb *KnowledgeBase) FindOne() error {
	collection := initializers.GetDb().C("KnowledgeBase")
	return collection.Find(bson.D{}).One(kb)
}

func (kb *KnowledgeBase) GetDataInput() []*DataInput {
	ret := []*DataInput{}
	for i := range kb.Objects {
		for j := range kb.Objects[i].Attributes {
			a := &kb.Objects[i].Attributes[j]
			if a.KbHistory == nil && a.KbAttribute.isSource(KBSource(User)) && !a.Validity() {
				di := DataInput{Id: a.Id, Name: a.KbObject.Name + "." + a.KbAttribute.Name, Atype: a.KbAttribute.AType, Options: a.KbAttribute.Options}
				ret = append(ret, &di)
			}
		}
	}
	return ret
}

func (kb *KnowledgeBase) FindAttributeObjectByName(name string) *KBAttributeObject {
	return kb.IdxAttributeObjects[name]
}
