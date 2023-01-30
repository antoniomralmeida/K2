package kb

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/kamva/mgm/v3"

	"github.com/antoniomralmeida/k2/cmd/k2/ebnf"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (kb *KnowledgeBased) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i := range attrs {
		attrs[i].ID = primitive.NewObjectID()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	inits.Log(c.Persist(), inits.Fatal)
}

func (kb *KnowledgeBased) Pause() {
	kb.halt = true
}

func (kb *KnowledgeBased) Resumo() {
	kb.halt = false
}

func (kb *KnowledgeBased) CopyClass(name string, copy *KBClass) *KBClass {
	if copy == nil {
		inits.Log(errors.New("Invalid class!"), inits.Error)
	}
	class := KBClass{}
	class.Name = name
	class.Attributes = copy.Attributes
	for i := range class.Attributes {
		class.Attributes[i].ID = primitive.NewObjectID()
	}
	err := class.Persist()
	if err == nil {
		kb.Classes = append(kb.Classes, class)
		kb.IdxClasses[class.ID] = &class
		return &class
	} else {
		inits.Log(err, inits.Error)
		return nil
	}
}

func (kb *KnowledgeBased) NewSimpleClass(name string, parent *KBClass) *KBClass {
	class := KBClass{}
	class.Name = name
	if parent != nil {
		class.ParentID = parent.ID
		class.ParentClass = parent
	}
	err := class.Persist()
	if err == nil {
		kb.Classes = append(kb.Classes, class)
		kb.IdxClasses[class.ID] = &class
		return &class
	} else {
		inits.Log(err, inits.Error)
		return nil
	}
}

func (kb *KnowledgeBased) NewClass(newclass_json string) *KBClass {
	class := KBClass{}
	err := json.Unmarshal([]byte(newclass_json), &class)
	if err != nil {
		inits.Log(err, inits.Info)
		return nil
	}
	if class.Parent != "" {
		p := kb.FindClassByName(class.Parent, true)
		if p == nil {
			inits.Log("Class not found "+class.Parent, inits.Info)
			return nil
		}
		class.ParentID = p.ID
		class.ParentClass = p
	}
	for i := range class.Attributes {
		class.Attributes[i].ID = primitive.NewObjectID()
		for _, x := range class.Attributes[i].Sources {
			class.Attributes[i].SourcesID = append(class.Attributes[i].SourcesID, KBSourceStr[x])
		}
		class.Attributes[i].SimulationID = KBSimulationStr[class.Attributes[i].Simulation]
	}
	err = class.Persist()
	if err == nil {
		kb.Classes = append(kb.Classes, class)
		kb.IdxClasses[class.ID] = &class
		return &class
	} else {
		inits.Log(err, inits.Error)
		return nil
	}
}

func (kb *KnowledgeBased) UpdateClass(c *KBClass) {
	for i := range c.Attributes {
		if c.Attributes[i].ID.IsZero() {
			c.Attributes[i].ID = primitive.NewObjectID()
		}
	}
	inits.Log(c.Persist(), inits.Fatal)
}

func (kb *KnowledgeBased) NewWorkspace(name string, image string) *KBWorkspace {
	copy, err := lib.LoadImage(image)
	if err != nil {
		inits.Log(err, inits.Error)
		return nil
	}
	w := KBWorkspace{Workspace: name, BackgroundImage: copy}
	err = w.Persist()
	if err == nil {
		kb.Workspaces = append(kb.Workspaces, w)
		return &w
	} else {
		inits.Log(err, inits.Fatal)
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

func (kb *KnowledgeBased) NewSimpleObject(name string, class *KBClass) *KBObject {
	o := KBObject{Name: name, Class: class.ID, Bkclass: class}
	for _, x := range kb.FindAttributes(class) {
		n := KBAttributeObject{Attribute: x.ID, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	inits.Log(o.Persist(), inits.Fatal)
	kb.IdxObjects[name] = &o
	return &o
}

func (kb *KnowledgeBased) NewObject(class string, name string) *KBObject {
	p := kb.FindClassByName(class, true)
	if p == nil {
		inits.Log("Class not found "+class, inits.Error)
		return nil
	}
	o := KBObject{Name: name, Class: p.ID, Bkclass: p}
	for _, x := range kb.FindAttributes(p) {
		n := KBAttributeObject{Attribute: x.ID, KbAttribute: x, KbObject: &o}
		o.Attributes = append(o.Attributes, n)
		kb.IdxAttributeObjects[n.getFullName()] = &n
	}
	inits.Log(o.Persist(), inits.Fatal)
	kb.IdxObjects[name] = &o
	return &o
}

func (kb *KnowledgeBased) LinkObjects(ws *KBWorkspace, obj *KBObject, left int, top int) {
	ows := KBObjectWS{Object: obj.ID, Left: left, Top: top, KBObject: obj}
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
		inits.Log(err, inits.Error)
		return nil
	}
	return kb.IdxClasses[ret.ID]
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

func (kb *KnowledgeBased) NewRule(rule string, priority byte, interval int) *KBRule {
	_, bin, err := kb.ParsingCommand(rule)
	if inits.Log(err, inits.Info) != nil {
		return nil
	}
	r := KBRule{Rule: rule, Priority: priority, ExecutionInterval: interval}
	inits.Log(r.Persist(), inits.Fatal)
	kb.linkerRule(&r, bin)
	kb.Rules = append(kb.Rules, r)
	return &r
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

func (kb *KnowledgeBased) GetWorkspaces() string {
	ret := []models.Workspace{}
	for _, w := range kb.Workspaces {
		ret = append(ret, models.Workspace{Workspace: w.Workspace, BackgroundImage: w.BackgroundImage})
	}
	json, err := json.Marshal(ret)
	inits.Log(err, inits.Error)
	return string(json)
}

func (kb *KnowledgeBased) RunStackRules() error {
	inits.Log("RunStackRules...", inits.Info)
	if len(kb.stack) > 0 {
		kb.mutex.Lock()
		localstack := kb.stack
		kb.mutex.Unlock()
		mark := len(localstack) - 1
		sort.Slice(localstack, func(i, j int) bool {
			return (localstack[i].Priority > localstack[j].Priority) || (localstack[i].Priority == localstack[j].Priority && localstack[j].lastexecution.Unix() > localstack[i].lastexecution.Unix())
		})

		runtaks := make(map[primitive.ObjectID]*KBRule) //run the rule once
		for _, r := range localstack {
			if runtaks[r.ID] == nil {
				r.Run()
				runtaks[r.ID] = r
			}
		}
		kb.mutex.Lock()
		kb.stack = kb.stack[mark:]
		kb.mutex.Unlock()
	}
	for i := range kb.Rules {
		if kb.Rules[i].ExecutionInterval != 0 && time.Now().After(kb.Rules[i].lastexecution.Add(time.Duration(kb.Rules[i].ExecutionInterval)*time.Millisecond)) {
			kb.mutex.Lock()
			kb.stack = append(kb.stack, &kb.Rules[i])
			kb.mutex.Unlock()
		}
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

func (kb *KnowledgeBased) ParsingCommand(cmd string) ([]*ebnf.Token, []*BIN, error) {
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
	var stack []*ebnf.Token
	var opts []*ebnf.Token
	var bin []*BIN
	for _, x := range tokens {
		var ok = false
		opts = kb.ebnf.FindOptions(pt, &stack, 0)
		for _, y := range opts {
			//fmt.Println(x, y)
			if (y.GetToken() == x) ||
				(y.GetTokentype() == ebnf.DynamicReference && len(x) == 1) ||
				((y.GetTokentype() == ebnf.Object || y.GetTokentype() == ebnf.Class || y.GetTokentype() == ebnf.Attribute || y.GetTokentype() == ebnf.Constant || y.GetTokentype() == ebnf.Reference) && unicode.IsUpper(rune(x[0]))) ||
				(y.GetTokentype() == ebnf.Text && (rune(x[0]) == '\'' || rune(x[0]) == '"') ||
					(y.GetTokentype() == ebnf.Constant && lib.IsNumber(x))) {
				if y.GetTokentype() == ebnf.Class {
					if kb.FindClassByName(x, false) != nil {
						ok = true
					}
				} else if y.GetTokentype() == ebnf.Object {
					if kb.FindObjectByName(x) != nil {
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
		if code.tokentype == ebnf.Literal && code.literalbin == models.B_null {
			inits.Log("Literal not found!", inits.Fatal)
		}
		bin = append(bin, &code)
	}
	for _, y := range pt.GetNexts() {
		if y.GetToken() == "." && y.GetTokentype() == ebnf.Control {
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

	dr := make(map[string]*KBClass)
	consequent := -1
	for j, x := range bin {
		switch x.literalbin {
		case models.B_initially:
			kb.mutex.Lock()
			kb.stack = append(kb.stack, r)
			kb.mutex.Unlock()
		case models.B_then:
			consequent = j
			r.consequent = j + 1
		}
		switch x.GetTokentype() {
		case ebnf.Workspace:
			if bin[j].workspace == nil {
				bin[j].workspace = kb.FindWorkspaceByName(r.bin[j].token)
			}
		case ebnf.Object:
			if len(bin[j].objects) == 0 {
				obj := kb.FindObjectByName(r.bin[j].token)
				bin[j].objects = append(bin[j].objects, obj)
			}
		case ebnf.Class:
			if bin[j].class == nil {
				c := kb.FindClassByName(x.GetToken(), true)
				bin[j].class = c
				objs := []KBObject{}
				inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Error)
				for _, y := range objs {
					bin[j].objects = append(bin[j].objects, kb.IdxObjects[y.Name])
				}
			}
		case ebnf.Attribute:
			ref := -1
			if bin[j+1].literalbin == models.B_of {
				ref = j + 2
			} else {
				for z := j - 1; z >= 0; z-- {
					if bin[z].GetTokentype() == ebnf.Object || bin[z].GetTokentype() == ebnf.Class {
						ref = z
						break
					}
				}
			}
			if ref != -1 {
				if bin[ref].GetTokentype() == ebnf.Object {
					if len(bin[j].objects) == 0 {
						obj := kb.FindObjectByName(r.bin[j].token)
						bin[j].objects = append(bin[j].objects, obj)
						bin[j].class = obj.Bkclass
					}
					bin[j].attribute = kb.FindAttribute(bin[ref].class, x.GetToken())
					if len(bin[j].objects) > 0 {
						atro := kb.FindAttributeObject(bin[ref].objects[0], x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				} else if bin[ref].GetTokentype() == ebnf.Class {
					c := bin[ref].class
					if c == nil {
						c = kb.FindClassByName(x.GetToken(), true)
						bin[ref].class = c
					}
					bin[j].class = c
					bin[j].attribute = kb.FindAttribute(c, x.GetToken())
					objs := []KBObject{}
					inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Fatal)
					for _, y := range objs {
						obj := kb.IdxObjects[y.Name]
						bin[j].objects = append(bin[j].objects, obj)
						atro := kb.FindAttributeObject(obj, x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				} else if bin[ref].GetTokentype() == ebnf.DynamicReference {
					c := bin[ref].class
					if c == nil {
						c = dr[bin[ref].token]
						bin[ref].class = c
					}
					if c == nil {
						return inits.Log("Attribute class not found in KB! "+x.GetToken(), inits.Error)
					}
					bin[j].attribute = kb.FindAttribute(c, x.GetToken())
					objs := []KBObject{}
					inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Fatal)
					for _, y := range objs {
						obj := kb.IdxObjects[y.Name]
						bin[j].objects = append(bin[j].objects, obj)
						atro := kb.FindAttributeObject(obj, x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				}
			} else {
				return inits.Log("Attribute not found in KB! "+x.GetToken(), inits.Error)
			}
		case ebnf.DynamicReference:
			{
				if consequent == -1 {
					for z := j - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == ebnf.Object || bin[z].GetTokentype() == ebnf.Class {
							bin[j].class = bin[z].class
							bin[j].objects = bin[z].objects
							dr[x.token] = bin[j].class
							break
						}
					}
				} else {
					for z := consequent - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == ebnf.DynamicReference && bin[z].GetToken() == x.GetToken() {
							bin[j].objects = bin[z].objects
							bin[j].class = bin[z].class
							dr[x.token] = bin[j].class
							break
						}
					}
				}
			}

		case ebnf.Constant:
			{
				if !lib.IsNumber(x.GetToken()) {
					ok := false
					for z := j - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == ebnf.Attribute {
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
	kb.mutex.Lock()
	r.bin = bin
	kb.mutex.Unlock()
	return nil
}
