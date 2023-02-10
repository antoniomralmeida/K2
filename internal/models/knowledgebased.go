package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/kamva/mgm/v3"
	"github.com/madflojo/tasks"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	_kb_current *KnowledgeBased
	scheduler   *tasks.Scheduler
)

type KnowledgeBased struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string        `bson:"name"`
	Classes          []KBClass     `bson:"-"`
	Rules            []KBRule      `bson:"-"`
	Workspaces       []KBWorkspace `bson:"-"`
	Objects          []KBObject    `bson:"-"`
	ebnf             *EBNF         `bson:"-"`
}

func (kb *KnowledgeBased) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i := range attrs {
		attrs[i].ID = primitive.NewObjectID()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	inits.Log(c.Persist(), inits.Fatal)
}

func KBPause() {
	scheduler.Lock()
}

func KBResume() {
	scheduler.Unlock()
}

func KBStop() {
	scheduler.Stop()
}

func KBCopyClass(name string, copy *KBClass) *KBClass {
	if _kb_current != nil {
		inits.Log(errors.New("Uninitialized KB!"), inits.Error)
		return nil
	}
	if copy == nil {
		inits.Log(errors.New("Invalid class!"), inits.Error)
		return nil
	}
	class := KBClass{}
	class.Name = name
	class.Attributes = copy.Attributes
	for i := range class.Attributes {
		class.Attributes[i].ID = primitive.NewObjectID()
	}
	err := class.Persist()
	if err == nil {
		_kb_current.Classes = append(_kb_current.Classes, class)
		//_kb.IdxClasses[class.ID] = &class
		return &class
	} else {
		inits.Log(err, inits.Error)
		return nil
	}
}

func KBNewSimpleClass(name string, parent *KBClass) *KBClass {
	if _kb_current != nil {
		inits.Log(errors.New("Uninitialized KB!"), inits.Error)
		return nil
	}
	class := KBClass{}
	class.Name = name
	if parent != nil {
		class.ParentID = parent.ID
		class.ParentClass = parent
	}
	err := class.Persist()
	if err == nil {
		_kb_current.Classes = append(_kb_current.Classes, class)
		//_kb.IdxClasses[class.ID] = &class
		return &class
	} else {
		inits.Log(err, inits.Error)
		return nil
	}
}

func (kb *KnowledgeBased) UpdateWorkspace(w *KBWorkspace) {
	inits.Log(w.Persist(), inits.Fatal)
}

func (kb *KnowledgeBased) FindWorkspaceByName(name string) *KBWorkspace {
	for i := range kb.Workspaces {
		if kb.Workspaces[i].Workspace == name {
			return &kb.Workspaces[i]
		}
	}
	inits.Log("Workspace not found!", inits.Error)
	return nil
}

func (kb *KnowledgeBased) LinkObjects(ws *KBWorkspace, obj *KBObject, left int, top int) {
	ows := KBObjectWS{Object: obj.ID, Left: left, Top: top, KBObject: obj}
	ws.Objects = append(ws.Objects, ows)
	kb.UpdateWorkspace(ws)
}

func (kb *KnowledgeBased) FindAttributeObject(obj *KBObject, attr string) *KBAttributeObject {
	for i := range obj.Attributes {
		if obj.Attributes[i].KbAttribute.Name == attr {
			return &obj.Attributes[i]
		}
	}
	return nil
}

func (kb *KnowledgeBased) AttributeObjectFactory(obj *KBObject, attr *KBAttribute) *KBAttributeObject {
	a := KBAttributeObject{Attribute: attr.ID}
	obj.Attributes = append(obj.Attributes, a)
	err := obj.Persist()
	if err == nil {
		return &a
	} else {
		inits.Log(err, inits.Fatal)
		return nil
	}
}

func (kb *KnowledgeBased) UpdateKB(name string) error {
	kb.Name = name
	return kb.Persist()
}

func (obj *KnowledgeBased) Persist() error {
	return inits.Persist(obj)

}

func (obj *KnowledgeBased) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (kb *KnowledgeBased) FindOne() error {
	ret := mgm.Coll(kb).FindOne(mgm.Ctx(), bson.D{})
	ret.Decode(kb)
	return nil
}

func KBGetDataInput() []*DataInput {
	objs := []KBObject{}
	mgm.Coll(new(KBObject)).SimpleFind(&objs, bson.D{})

	ret := []*DataInput{}
	for i := range objs {
		for j := range objs[i].Attributes {
			a := &objs[i].Attributes[j]
			if a.KbAttribute.isSource(FromUser) && !a.Validity() {
				di := DataInput{Name: a.KbObject.Name + "." + a.KbAttribute.Name, Atype: a.KbAttribute.AType, Options: a.KbAttribute.Options}
				ret = append(ret, &di)
			}
		}
	}
	return ret
}

func KBFindAttributeObjectByName(key string) *KBAttributeObject {
	keys := strings.Split(key, ".")
	ao := new(KBObject)
	r := mgm.Coll(ao).FindOne(mgm.Ctx(), bson.D{{"name", keys[0]}, {"attribute.name", key[1]}})
	r.Decode(ao)
	return &ao.Attributes[0]
}

func KBGetWorkspaces() string {
	wks := []KBWorkspace{}

	mgm.Coll(new(KBWorkspace)).SimpleFind(&wks, bson.D{{}})
	ret := []WorkspaceInfo{}
	for _, w := range wks {
		ret = append(ret, WorkspaceInfo{Workspace: w.Workspace, BackgroundImage: w.BackgroundImage})
	}
	json, err := json.Marshal(ret)
	inits.Log(err, inits.Error)
	return string(json)
}

func KBGetWorkspacesFromObject(o *KBObject) (ret []*KBWorkspace) {
	//TODO: From mongoDB
	for i := range _kb_current.Workspaces {
		for j := range _kb_current.Workspaces[i].Objects {
			if _kb_current.Workspaces[i].Objects[j].KBObject == o {
				ret = append(ret, &_kb_current.Workspaces[i])
			}
		}
	}
	return
}

func (kb *KnowledgeBased) RunStackRules() error {
	inits.Log("RunStackRules...", inits.Info)
	for i := range kb.Rules {
		if kb.Rules[i].ExecutionInterval != 0 && time.Now().After(kb.Rules[i].Lastexecution.Add(time.Duration(kb.Rules[i].ExecutionInterval)*time.Millisecond)) {
			stack := KBStack{RuleID: kb.Rules[i].ID}
			stack.Persist()
		}
	}

	toRun := StacktoRun()
	for _, r := range toRun {
		r.Run()
	}

	return nil
}

func (kb *KnowledgeBased) RefreshRules() error {
	inits.Log("RefrehRules...", inits.Info)
	for i := range kb.Objects {
		if !kb.Objects[i].parsed {
			for j := range kb.Rules {
				for k := range kb.Rules[j].bkclasses {
					if kb.Rules[j].bkclasses[k] == kb.Objects[i].Bkclass {
						_, bin, err := kb.ParsingCommand(kb.Rules[j].Rule)
						if inits.Log(err, inits.Error) != nil {
							kb.linkerRule(&kb.Rules[j], bin)
						}
					}
				}
			}
			kb.Objects[i].parsed = true
		}
	}
	return nil
}

func (kb *KnowledgeBased) ParsingCommand(cmd string) ([]*Token, []*BIN, error) {
	cmd = strings.Replace(cmd, "\r\n", "", -1)
	cmd = strings.Replace(cmd, "\\n", "", -1)
	cmd = strings.Replace(cmd, "\t", " ", -1)
	for strings.Contains(cmd, "  ") {
		cmd = strings.Replace(cmd, "  ", " ", -1)
	}
	inits.Log("Parsing Prodution Rule: "+cmd, inits.Info)
	var inWord = false
	var inString = false
	var inNumber = false
	var start = 0
	var tokens []string
	const endline = '春'
	cmd = cmd + string(endline)
	for i, c := range cmd {
		switch {
		case c == '春' || c == ' ' || kb.ebnf.FindSymbols(string(c), true) != -1:
			if inNumber && c != '.' {
				tokens = append(tokens, cmd[start:i])
				inNumber = false
			} else if inString {
				if c == '"' || c == '\'' {
					tokens = append(tokens, cmd[start:i+1])
					inString = false
				}
			} else if inWord {
				tokens = append(tokens, cmd[start:i])
				inWord = false
			} else {
				if c == '"' || c == '\'' {
					start = i
					inString = true
				} else if c != ' ' && c != '.' && c != endline {
					tokens = append(tokens, string(c))
				}
			}
		case unicode.IsLower(c) && !inWord && !inString && !inNumber:
			start = i
			inWord = true
		case unicode.IsUpper(c) && !inWord && !inString && !inNumber:
			start = i
			inWord = true
		case unicode.IsNumber(c) && !inNumber && !inString && !inWord:
			start = i
			inNumber = true
		default:
		}
	}
	var pt = kb.ebnf.GetBase()
	var stack []*Token
	var opts []*Token
	var bin []*BIN
	for _, x := range tokens {
		var ok = false
		opts = kb.ebnf.FindOptions(pt, &stack, 0)
		for _, y := range opts {
			//fmt.Println(x, y)
			if (y.GetToken() == x) ||
				(y.GetTokentype() == DynamicReference && len(x) == 1) ||
				((y.GetTokentype() == Object || y.GetTokentype() == Class ||
					y.GetTokentype() == Attribute || y.GetTokentype() == Constant ||
					y.GetTokentype() == Reference) && unicode.IsUpper(rune(x[0]))) ||
				(y.GetTokentype() == Text && (rune(x[0]) == '\'' || rune(x[0]) == '"') ||
					(y.GetTokentype() == Constant && lib.IsNumber(x))) {
				if y.GetTokentype() == Class {
					if FindClassByName(x, false) != nil {
						ok = true
					}
				} else if y.GetTokentype() == Object {
					if FindObjectByName(x) != nil {
						ok = true
					}
				} else {
					ok = true
				}
				if ok {
					pt = y
					break

				}
			}
		}
		if !ok || len(opts) == 0 {
			str := "Compiller error in " + x + " when the expected was: "
			for _, y := range opts {
				str = str + "... " + y.GetToken()
			}
			return opts, nil, errors.New(str)
		}
		code := BIN{tokentype: pt.GetTokentype(), token: x}
		code.setTokenBin()
		if code.tokentype == Literal && code.literalbin == B_null {
			inits.Log("Literal not found!", inits.Fatal)
		}
		bin = append(bin, &code)
	}
	for _, y := range pt.GetNexts() {
		if y.GetToken() == "." && y.GetTokentype() == Control {
			inits.Log(", compilation successfully!", inits.Info)
			return nil, bin, nil
		}
	}
	opts = kb.ebnf.FindOptions(pt, &stack, 0)
	str := "Incomplete sentence when the expected was: "
	for _, y := range opts {
		str = str + "... " + y.GetToken()
	}
	return opts, nil, errors.New(str)
}

func (kb *KnowledgeBased) linkerRule(r *KBRule, bin []*BIN) error {
	// Find references of objects in KB
	inits.Log("Linking Prodution Rule: "+r.ID.Hex(), inits.Info)
	KBPause()

	dr := make(map[string]*KBClass)
	consequent := -1
	for j, x := range bin {
		switch x.literalbin {
		case B_initially:
			stack := KBStack{RuleID: r.ID}
			stack.Persist()
		case B_then:
			consequent = j
			r.consequent = j + 1
		}
		switch x.GetTokentype() {
		case Workspace:
			if bin[j].workspace == nil {
				bin[j].workspace = kb.FindWorkspaceByName(r.bin[j].token)
			}
		case Object:
			if len(bin[j].objects) == 0 {
				obj := FindObjectByName(r.bin[j].token)
				bin[j].objects = append(bin[j].objects, obj)
			}
		case Class:
			if bin[j].class == nil {
				c := FindClassByName(x.GetToken(), true)
				bin[j].class = c
				objs := []KBObject{}
				inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Error)
				for _, y := range objs {
					bin[j].objects = append(bin[j].objects, &y)
				}
			}
		case Attribute:
			ref := -1
			if bin[j+1].literalbin == B_of {
				ref = j + 2
			} else {
				for z := j - 1; z >= 0; z-- {
					if bin[z].GetTokentype() == Object || bin[z].GetTokentype() == Class {
						ref = z
						break
					}
				}
			}
			if ref != -1 {
				if bin[ref].GetTokentype() == Object {
					if len(bin[j].objects) == 0 {
						obj := FindObjectByName(r.bin[j].token)
						bin[j].objects = append(bin[j].objects, obj)
						bin[j].class = obj.Bkclass
					}
					bin[j].attribute = FindAttribute(bin[ref].class, x.GetToken())
					if len(bin[j].objects) > 0 {
						atro := kb.FindAttributeObject(bin[ref].objects[0], x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				} else if bin[ref].GetTokentype() == Class {
					c := bin[ref].class
					if c == nil {
						c = FindClassByName(x.GetToken(), true)
						bin[ref].class = c
					}
					bin[j].class = c
					bin[j].attribute = FindAttribute(c, x.GetToken())
					objs := []KBObject{}
					inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Fatal)
					for _, y := range objs {
						obj := &y
						bin[j].objects = append(bin[j].objects, obj)
						atro := kb.FindAttributeObject(obj, x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				} else if bin[ref].GetTokentype() == DynamicReference {
					c := bin[ref].class
					if c == nil {
						c = dr[bin[ref].token]
						bin[ref].class = c
					}
					if c == nil {
						return inits.Log("Attribute class not found in KB! "+x.GetToken(), inits.Error)
					}
					bin[j].attribute = FindAttribute(c, x.GetToken())
					objs := []KBObject{}
					inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Fatal)
					for _, y := range objs {
						obj := &y
						bin[j].objects = append(bin[j].objects, obj)
						atro := kb.FindAttributeObject(obj, x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				}
			} else {
				return inits.Log("Attribute not found in KB! "+x.GetToken(), inits.Error)
			}
		case DynamicReference:
			{
				if consequent == -1 {
					for z := j - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == Object || bin[z].GetTokentype() == Class {
							bin[j].class = bin[z].class
							bin[j].objects = bin[z].objects
							dr[x.token] = bin[j].class
							break
						}
					}
				} else {
					for z := consequent - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == DynamicReference && bin[z].GetToken() == x.GetToken() {
							bin[j].objects = bin[z].objects
							bin[j].class = bin[z].class
							dr[x.token] = bin[j].class
							break
						}
					}
				}
			}

		case Constant:
			{
				if !lib.IsNumber(x.GetToken()) {
					ok := false
					for z := j - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == Attribute {
							if bin[z].attribute != nil {
								for _, o := range bin[z].attribute.Options {
									if x.GetToken() == o {
										bin[j].token = "\"" + bin[j].token + "\""
										ok = true
										break
									}
								}
							}
						}
					}
					if !ok {
						return inits.Log("List option not found in KB! "+x.GetToken(), inits.Error)
					}
				}
			}
		}
		a := bin[j].attribute
		if a != nil {
			if consequent != -1 {
				a.addConsequentRules(r)
			} else {
				a.addAntecedentRules(r)
			}
		}
		cl := bin[j].class
		if cl != nil {
			r.addClass(cl)
		}
		for z := range bin[j].objects {
			bin[j].objects[z].parsed = true
		}
	}
	r.bin = bin
	KBResume()
	return nil
}

func KnowledgeBasedFacotory() *KnowledgeBased {
	kb := new(KnowledgeBased)
	kb.FindOne()
	if kb.Name == "" {
		kb.Name = "K2 KnowledgeBase System "
	}
	kb.Persist()
	return kb
}

func KBInit() {
	inits.Log("Init KB", inits.Info)
	_kb_current = KnowledgeBasedFacotory()

	ebnf := EBNF{}
	_kb_current.ebnf = &ebnf
	_kb_current.ebnf.ReadToken("./configs/k2.ebnf")

	FindAllClasses("_id", &_kb_current.Classes)

	_idxClasses := make(map[primitive.ObjectID]*KBClass)
	for _, c := range _kb_current.Classes {
		_idxClasses[c.ID] = &c
	}

	for j, c := range _kb_current.Classes {
		inits.Log("Prepare Class "+c.Name, inits.Info)
		if !c.ParentID.IsZero() {
			pc := _idxClasses[c.ParentID]
			if pc != nil {
				_kb_current.Classes[j].ParentClass = pc
			} else {
				inits.Log("Parent of Class "+c.Name+" not found!", inits.Fatal)
			}
		}
	}

	FindAllObjects(bson.M{}, "name", &_kb_current.Objects)
	for j, o := range _kb_current.Objects {
		//_kb.IdxObjects[o.Name] = &_kb.Objects[j]
		c := _idxClasses[o.Class]
		if c != nil {
			_kb_current.Objects[j].Bkclass = c
			attrs := FindAttributes(c)
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].ID.Hex() < attrs[j].ID.Hex()
			})
			for k, x := range o.Attributes {
				_kb_current.Objects[j].Attributes[k].KbObject = &_kb_current.Objects[j]
				//kb.Objects[j].Attributes[k].Kb = kb
				for l, y := range attrs {
					if y.ID == x.Attribute {
						_kb_current.Objects[j].Attributes[k].KbAttribute = attrs[l]
						break
					}
					if y.ID.Hex() > x.Attribute.Hex() {
						break
					}
				}
				if _kb_current.Objects[j].Attributes[k].KbAttribute == nil {
					inits.Log("Attribute not found "+x.Attribute.Hex(), inits.Fatal)
				}
				//_kb.IdxAttributeObjects[o.Name+"."+_kb.Objects[j].Attributes[k].KbAttribute.Name] = &_kb.Objects[j].Attributes[k]

				//Obter ultimo valor
				h := KBHistory{}
				err := h.FindLast(bson.D{{Key: "attribute_id", Value: x.ID}})
				if err != nil {
					if err.Error() != "not found" {
						inits.Log(err, inits.Fatal)
					}
				} else {
					_kb_current.Objects[j].Attributes[k].KbHistory = &h
				}
				_kb_current.Objects[j].Attributes[k].Validity()
			}
		} else {
			inits.Log("Class of object "+o.Name+" not found!", inits.Fatal)
		}
	}

	FindAllWorkspaces("name")

	FindAllRules("_id")

	for i := range _kb_current.Rules {
		_, bin, err := _kb_current.ParsingCommand(_kb_current.Rules[i].Rule)
		inits.Log(err, inits.Fatal)
		_kb_current.linkerRule(&_kb_current.Rules[i], bin)
	}
}

func FindAllWorkspaces(sort string) error {
	collection := mgm.Coll(new(KBWorkspace))
	idx := collection.Indexes()
	ret, err := idx.List(mgm.Ctx())
	inits.Log(err, inits.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	inits.Log(err, inits.Fatal)
	if len(results) == 1 {
		inits.CreateUniqueIndex(collection, "workspace")
	}
	cursor, err := collection.Find(mgm.Ctx(), bson.D{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), &_kb_current.Workspaces)
	return err
}

func FindAllClasses(sort string, cs *[]KBClass) error {
	collection := mgm.Coll(new(KBClass))
	idx := collection.Indexes()
	ret, err := idx.List(context.TODO())
	inits.Log(err, inits.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	inits.Log(err, inits.Fatal)
	if len(results) == 1 {
		inits.CreateUniqueIndex(collection, "name")
	}
	cursor, err := collection.Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), cs)
	return err
}

func FindAllObjects(filter bson.M, sort string, os *[]KBObject) error {
	collection := mgm.Coll(new(KBObject))
	idx := collection.Indexes()
	ret, err := idx.List(mgm.Ctx())
	inits.Log(err, inits.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	inits.Log(err, inits.Fatal)
	if len(results) == 1 {
		inits.CreateUniqueIndex(collection, "name")
	}
	cursor, err := collection.Find(mgm.Ctx(), filter, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), os)
	return err
}

func FindAllRules(sort string) error {
	collection := mgm.Coll(new(KBRule))
	cursor, err := collection.Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), &_kb_current.Rules)
	return err
}

func KBRun(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start the Scheduler
	scheduler = tasks.New()
	defer scheduler.Stop()

	// Add tasks
	_, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(2 * time.Second),

		TaskFunc: func() error {
			go _kb_current.RunStackRules()
			return nil
		},
	})
	inits.Log(err, inits.Fatal)
	_, err = scheduler.Add(&tasks.Task{
		Interval: time.Duration(60 * time.Second),
		TaskFunc: func() error {
			go _kb_current.RefreshRules()
			return nil
		},
	})
	inits.Log(err, inits.Fatal)

	inits.Log("K2 KB System started!", inits.Info)
	if runtime.GOOS == "windows" {
		fmt.Println("K2 KB System started! Press ESC to shutdown")
	}
	for {
		if lib.KeyPress() == 27 {
			fmt.Printf("Shutdown...")
			KBStop()
			wg.Done()
			os.Exit(0)
		}
	}

}
