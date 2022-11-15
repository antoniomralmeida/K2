package kb

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var GKB *KnowledgeBased

func Init() {
	initializers.Log("Init KB", initializers.Info)
	k := KnowledgeBased{}
	GKB = &k
	GKB.FindOne()
	if GKB.Name == "" {
		GKB.Name = "K2 System KB"
	}
	GKB.Persist()
	GKB.IdxClasses = make(map[primitive.ObjectID]*KBClass)
	GKB.IdxObjects = make(map[string]*KBObject)
	GKB.IdxAttributeObjects = make(map[string]*KBAttributeObject)

	ebnf := ebnf.EBNF{}
	GKB.ebnf = &ebnf
	GKB.ebnf.ReadToken("./k2web/pub/ebnf/k2.ebnf")

	FindAllClasses("_id", &GKB.Classes)
	for j := range GKB.Classes {
		GKB.IdxClasses[GKB.Classes[j].Id] = &GKB.Classes[j]
	}

	for j, c := range GKB.Classes {
		initializers.Log("Prepare Class "+c.Name, initializers.Info)
		if !c.ParentID.IsZero() {
			pc := GKB.IdxClasses[c.ParentID]
			if pc != nil {
				GKB.Classes[j].ParentClass = pc
			} else {
				initializers.Log("Parent of Class "+c.Name+" not found!", initializers.Fatal)
			}
		}
	}

	FindAllObjects(bson.M{}, "name", &GKB.Objects)
	for j, o := range GKB.Objects {
		GKB.IdxObjects[o.Name] = &GKB.Objects[j]
		c := GKB.IdxClasses[o.Class]
		if c != nil {
			GKB.Objects[j].Bkclass = c
			attrs := GKB.FindAttributes(c)
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].Id.Hex() < attrs[j].Id.Hex()
			})
			for k, x := range o.Attributes {
				GKB.Objects[j].Attributes[k].KbObject = &GKB.Objects[j]
				//kb.Objects[j].Attributes[k].Kb = kb
				for l, y := range attrs {
					if y.Id == x.Attribute {
						GKB.Objects[j].Attributes[k].KbAttribute = attrs[l]
						break
					}
					if y.Id.Hex() > x.Attribute.Hex() {
						break
					}
				}
				if GKB.Objects[j].Attributes[k].KbAttribute == nil {
					initializers.Log("Attribute not found "+x.Attribute.Hex(), initializers.Fatal)
				}
				GKB.IdxAttributeObjects[o.Name+"."+GKB.Objects[j].Attributes[k].KbAttribute.Name] = &GKB.Objects[j].Attributes[k]

				//Obter ultimo valor
				h := KBHistory{}
				err := h.FindLast(bson.D{{Key: "attribute_id", Value: x.Id}})
				if err != nil {
					if err.Error() != "not found" {
						initializers.Log(err, initializers.Fatal)
					}
				} else {
					GKB.Objects[j].Attributes[k].KbHistory = &h
				}
				GKB.Objects[j].Attributes[k].Validity()
			}
		} else {
			initializers.Log("Class of object "+o.Name+" not found!", initializers.Fatal)
		}
	}

	FindAllWorkspaces("name")

	FindAllRules("_id")

	for i := range GKB.Rules {
		_, bin, err := GKB.ParsingCommand(GKB.Rules[i].Rule)
		initializers.Log(err, initializers.Fatal)
		GKB.linkerRule(&GKB.Rules[i], bin)
	}
}

func (kb *KnowledgeBased) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i := range attrs {
		attrs[i].Id = primitive.NewObjectID()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	initializers.Log(c.Persist(), initializers.Fatal)
}

func (kb *KnowledgeBased) NewClass(newclass_json string) *KBClass {
	class := KBClass{}
	err := json.Unmarshal([]byte(newclass_json), &class)
	if err != nil {
		initializers.Log(err, initializers.Info)
		return nil
	}
	if class.Parent != "" {
		p := kb.FindClassByName(class.Parent, true)
		if p == nil {
			initializers.Log("Class not found "+class.Parent, initializers.Info)
			return nil
		}
		class.ParentID = p.Id
		class.ParentClass = p
	}
	for i := range class.Attributes {
		class.Attributes[i].Id = primitive.NewObjectID()
		for _, x := range class.Attributes[i].Sources {
			class.Attributes[i].SourcesID = append(class.Attributes[i].SourcesID, KBSourceStr[x])
		}
		class.Attributes[i].SimulationID = KBSimulationStr[class.Attributes[i].Simulation]
	}
	err = class.Persist()
	if err == nil {
		kb.Classes = append(kb.Classes, class)
		kb.IdxClasses[class.Id] = &class
		return &class
	} else {
		initializers.Log(err, initializers.Error)
		return nil
	}
}

func (kb *KnowledgeBased) UpdateClass(c *KBClass) {
	for i := range c.Attributes {
		if c.Attributes[i].Id.IsZero() {
			c.Attributes[i].Id = primitive.NewObjectID()
		}
	}
	initializers.Log(c.Persist(), initializers.Fatal)
}

func (kb *KnowledgeBased) NewWorkspace(name string, icone string) *KBWorkspace {
	w := KBWorkspace{Workspace: name, BackgroundImage: icone}
	err := w.Persist()
	if err == nil {
		kb.Workspaces = append(kb.Workspaces, w)
		return &w
	} else {
		initializers.Log(err, initializers.Error)
		return nil
	}
}

func (kb *KnowledgeBased) UpdateWorkspace(w *KBWorkspace) {
	initializers.Log(w.Persist(), initializers.Fatal)
}

func (kb *KnowledgeBased) FindWorkspaceByName(name string) *KBWorkspace {
	for i := range kb.Workspaces {
		if kb.Workspaces[i].Workspace == name {
			return &kb.Workspaces[i]
		}
	}
	initializers.Log("Workspace not found!", initializers.Error)
	return nil
}

func (kb *KnowledgeBased) NewObject(class string, name string) *KBObject {
	p := kb.FindClassByName(class, true)
	if p == nil {
		initializers.Log("Class not found "+class, initializers.Error)
		return nil
	}
	o := KBObject{Name: name, Class: p.Id, Bkclass: p}
	for _, x := range kb.FindAttributes(p) {
		n := KBAttributeObject{Id: primitive.NewObjectID(), Attribute: x.Id, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	initializers.Log(o.Persist(), initializers.Fatal)
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
	err := ret.FindOne(bson.D{{Key: "name", Value: nm}})
	if err != nil && mandatory {
		initializers.Log(err, initializers.Error)
		return nil
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
	a := KBAttributeObject{Attribute: attr.Id, Id: primitive.NewObjectID()}
	obj.Attributes = append(obj.Attributes, a)
	err := obj.Persist()
	if err == nil {
		return &a
	} else {
		initializers.Log(err, initializers.Error)
		return nil
	}
}

func (kb *KnowledgeBased) NewRule(rule string, priority byte, interval int) *KBRule {
	_, bin, err := kb.ParsingCommand(rule)
	if initializers.Log(err, initializers.Info) != nil {
		return nil
	}
	r := KBRule{Rule: rule, Priority: priority, ExecutionInterval: interval}
	initializers.Log(r.Persist(), initializers.Fatal)
	kb.linkerRule(&r, bin)
	kb.Rules = append(kb.Rules, r)
	return &r
}
func (kb *KnowledgeBased) UpdateKB(name string) error {
	kb.Name = name
	return kb.Persist()
}

func (kb *KnowledgeBased) PrintEBNF() {
	fmt.Println(kb.ebnf.String())
}

func (kb *KnowledgeBased) Persist() error {
	ctx, collection := initializers.GetCollection("KnowledgeBased")
	if kb.Id.IsZero() {
		kb.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, kb)
		return err
	} else {

		_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: kb.Id}}, kb)
		return err
	}
}

func (kb *KnowledgeBased) FindOne() error {
	ctx, collection := initializers.GetCollection("KnowledgeBased")
	ret := collection.FindOne(ctx, bson.D{})
	ret.Decode(kb)
	return nil
}

func (kb *KnowledgeBased) GetDataInput() []*DataInput {
	ret := []*DataInput{}
	for i := range kb.Objects {
		for j := range kb.Objects[i].Attributes {
			a := &kb.Objects[i].Attributes[j]
			if a.KbAttribute.isSource(User) && !a.Validity() {
				di := DataInput{Name: a.KbObject.Name + "." + a.KbAttribute.Name, Atype: a.KbAttribute.AType, Options: a.KbAttribute.Options}
				ret = append(ret, &di)
			}
		}
	}
	return ret
}

func (kb *KnowledgeBased) FindAttributeObjectByName(name string) *KBAttributeObject {
	return kb.IdxAttributeObjects[name]
}
