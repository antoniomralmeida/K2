package kb

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"gopkg.in/mgo.v2/bson"
)

func (kb *KnowledgeBased) Init() {
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
		if c.ParentID != "" {
			pc := kb.IdxClasses[c.ParentID]
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

func (kb *KnowledgeBased) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i := range attrs {
		attrs[i].Id = bson.NewObjectId()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	lib.LogFatal(c.Persist())
}

func (kb *KnowledgeBased) NewClass(newclass_json string) *KBClass {
	class := KBClass{}
	err := json.Unmarshal([]byte(newclass_json), &class)
	if err != nil {
		log.Println(err)
		return nil
	}
	if class.Parent != "" {
		p := kb.FindClassByName(class.Parent, true)
		if p == nil {
			log.Println("Class not found ", class.Parent)
			return nil
		}
		class.ParentID = p.Id
		class.ParentClass = p
	}
	for i := range class.Attributes {
		class.Attributes[i].Id = bson.NewObjectId()
		for _, x := range class.Attributes[i].Sources {
			class.Attributes[i].SourcesID = append(class.Attributes[i].SourcesID, KBSourceStr[x])
		}
	}
	err = class.Persist()
	if err != nil {
		log.Fatal(err)
	}
	kb.Classes = append(kb.Classes, class)
	kb.IdxClasses[class.Id] = &class
	return &class
}

func (kb *KnowledgeBased) UpdateClass(c *KBClass) {
	for i := range c.Attributes {
		if c.Attributes[i].Id == "" {
			c.Attributes[i].Id = bson.NewObjectId()
		}
	}
	lib.LogFatal(c.Persist())
}

func (kb *KnowledgeBased) NewWorkspace(name string, icone string) *KBWorkspace {
	w := KBWorkspace{Workspace: name, BackgroundImage: icone}
	log.Fatal(w.Persist())
	kb.Workspaces = append(kb.Workspaces, w)
	return &w
}

func (kb *KnowledgeBased) UpdateWorkspace(w *KBWorkspace) {
	lib.LogFatal(w.Persist())
}

func (kb *KnowledgeBased) FindWorkspaceByName(name string) *KBWorkspace {
	for i := range kb.Workspaces {
		if kb.Workspaces[i].Workspace == name {
			return &kb.Workspaces[i]
		}
	}
	log.Fatal("Workspace not found!")
	return nil
}

func (kb *KnowledgeBased) NewObject(class string, name string) *KBObject {
	p := kb.FindClassByName(class, true)
	if p == nil {
		log.Println("Class not found ", class)
		return nil
	}
	o := KBObject{Name: name, Class: p.Id, Bkclass: p}
	for _, x := range kb.FindAttributes(p) {
		n := KBAttributeObject{Id: bson.NewObjectId(), Attribute: x.Id, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	lib.LogFatal(o.Persist())
	kb.IdxObjects[name] = &o
	return &o
}

func (kb *KnowledgeBased) LinkObjects(ws *KBWorkspace, obj *KBObject, left int, top int) {
	ows := KBObjectWS{Object: obj.Id, Left: left, Top: top, KBObject: obj}
	ws.Objects = append(ws.Objects, ows)
	kb.UpdateWorkspace(ws)
}

func (kb *KnowledgeBased) FindObjectByName(name string) *KBObject {
	return kb.IdxObjects[name]
}

func (kb *KnowledgeBased) FindClassByName(nm string, mandatory bool) *KBClass {
	var ret KBClass
	err := ret.FindOne(bson.D{{"name", nm}})
	if err != nil && mandatory {
		log.Fatal(err)
	}
	return kb.IdxClasses[ret.Id]
}

func (kb *KnowledgeBased) FindAttribute(c *KBClass, name string) *KBAttribute {
	attrs := kb.FindAttributes(c)
	for i, x := range attrs {
		if x.Name == name {
			return attrs[i]
		}
	}
	return nil
}

func (kb *KnowledgeBased) FindAttributes(c *KBClass) []*KBAttribute {
	var ret []*KBAttribute
	if c != nil {
		if c.ParentClass != nil {
			ret = append(ret, kb.FindAttributes(c.ParentClass)...)
		}
		for i := range c.Attributes {
			ret = append(ret, &c.Attributes[i])
		}
	}
	return ret
}

func (kb *KnowledgeBased) FindAttributeObject(obj *KBObject, attr string) *KBAttributeObject {
	for i := range obj.Attributes {
		if obj.Attributes[i].KbAttribute.Name == attr {
			return &obj.Attributes[i]
		}
	}
	return nil
}

func (kb *KnowledgeBased) NewAttributeObject(obj *KBObject, attr *KBAttribute) *KBAttributeObject {
	a := KBAttributeObject{Attribute: attr.Id, Id: bson.NewObjectId()}
	obj.Attributes = append(obj.Attributes, a)
	log.Fatal(obj.Persist())
	return &a
}

func (kb *KnowledgeBased) NewRule(rule string, priority byte, interval int) *KBRule {
	_, bin, err := kb.ParsingCommand(rule)
	lib.LogFatal(err)
	r := KBRule{Rule: rule, Priority: priority, ExecutionInterval: interval}
	lib.LogFatal(r.Persist())
	kb.linkerRule(&r, bin)
	kb.Rules = append(kb.Rules, r)
	return &r
}
func (kb *KnowledgeBased) UpdateKB(name string, iotapi string) error {
	kb.Name = name
	kb.IOTApi = iotapi
	return kb.Persist()
}

func (kb *KnowledgeBased) PrintEBNF() {
	kb.ebnf.PrintEBNF()
}

func (kb *KnowledgeBased) Persist() error {
	collection := initializers.GetDb().C("KnowledgeBased")
	if kb.Id == "" {
		kb.Id = bson.NewObjectId()
		return collection.Insert(kb)
	} else {
		return collection.UpdateId(kb.Id, kb)
	}
}

func (kb *KnowledgeBased) FindOne() error {
	collection := initializers.GetDb().C("KnowledgeBased")
	return collection.Find(bson.D{}).One(kb)
}

func (kb *KnowledgeBased) GetDataInput() []*DataInput {
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

func (kb *KnowledgeBased) FindAttributeObjectByName(name string) *KBAttributeObject {
	return kb.IdxAttributeObjects[name]
}
