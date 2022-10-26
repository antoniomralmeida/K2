package kb

import (
	"errors"
	"log"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/lib"
	"gopkg.in/mgo.v2/bson"
)

func (kb *KnowledgeBase) Init() {
	kb.IdxClasses = make(map[bson.ObjectId]KBClassPt)
	kb.IdxObjects = make(map[string]*KBObject)
	TokenBinStr = map[string]TokenBin{
		"":                b_null,
		"(":               b_open_par,
		")":               b_close_par,
		"=":               b_equal_sym,
		"activate":        b_activate,
		"and":             b_and,
		"any":             b_any,
		"change":          b_change,
		"conclude":        b_conclude,
		"create":          b_create,
		"deactivate":      b_deactivate,
		"delete":          b_delete,
		"different":       b_different,
		"equal":           b_equal,
		"focus":           b_focus,
		"for":             b_for,
		"greater":         b_greater,
		"halt":            b_halt,
		"hide":            b_hide,
		"if":              b_if,
		"inform":          b_inform,
		"initially":       b_initially,
		"insert":          b_insert,
		"invoke":          b_invoke,
		"is":              b_is,
		"less":            b_less,
		"move":            b_move,
		"of":              b_of,
		"operator":        b_operator,
		"or":              b_or,
		"remove":          b_remove,
		"rotate":          b_rotate,
		"set":             b_set,
		"show":            b_show,
		"start":           b_start,
		"than":            b_than,
		"that":            b_than,
		"the":             b_the,
		"then":            b_then,
		"to":              b_to,
		"transfer":        b_transfer,
		"unconditionally": b_unconditionally,
		"when":            b_when,
		"whenever":        b_whenever}

}

func (kb *KnowledgeBase) Scan() error {
	lib.Log("Scaning...")
	if len(kb.stack) > 0 {
		localstack := kb.stack
		mark := len(localstack) - 1
		sort.Slice(localstack, func(i, j int) bool {
			return (localstack[i].Priority > localstack[j].Priority) || (localstack[i].Priority == localstack[j].Priority && localstack[j].lastexecution.After(localstack[i].lastexecution))
		})

		for len(localstack) > 0 {
			r := localstack[0]
			r.Run()
			localstack = localstack[1:]
		}
		kb.stack = kb.stack[mark:]
	}
	for i, _ := range kb.Rules {
		if kb.Rules[i].ExecutionInterval != 0 && time.Now().After(kb.Rules[i].lastexecution.Add(time.Duration(kb.Rules[i].ExecutionInterval)*time.Millisecond)) {
			kb.stack = append(kb.stack, &kb.Rules[i])
		}
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (kb *KnowledgeBase) ParsingCommand(cmd string) ([]*ebnf.Token, []*BIN, error) {
	cmd = strings.Replace(cmd, "\r\n", "", -1)
	cmd = strings.Replace(cmd, "\\n", "", -1)
	cmd = strings.Replace(cmd, "\t", " ", -1)
	for strings.Contains(cmd, "  ") {
		cmd = strings.Replace(cmd, "  ", " ", -1)
	}
	log.Println("Parsing Prodution Rule: ", cmd)
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
		bin = append(bin, &code)
	}
	for _, y := range pt.GetNexts() {
		if y.GetToken() == "." && y.GetTokentype() == ebnf.Control {
			log.Println(", compilation successfully!")
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
		switch x.GetTokentype() {
		case ebnf.Object:
			if bin[j].object == nil {
				bin[j].object = kb.FindObjectByName(x.GetToken())
			}
		case ebnf.Class:
			if bin[j].class == nil {
				bin[j].class = kb.FindClassByName(x.GetToken(), true)
			}
		case ebnf.Attribute:
			ref := -1
			if bin[j+1].typebin == b_of {
				ref = j + 3
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
					if bin[ref].object == nil {
						bin[ref].object = kb.FindObjectByName(bin[ref].GetToken())
					}
					bin[j].attribute = kb.FindAttributeObjectByClass(bin[ref].object, x.GetToken())
					break
				} else if bin[ref].GetTokentype() == ebnf.Class {
					if bin[ref].class == nil {
						bin[j].class = kb.FindClassByName(x.GetToken(), true)
					}
					bin[j].attribute = kb.FindAttribute(bin[ref].class, x.GetToken())
					break
				}
			} else {
				log.Fatal("Attribute not found in KB! ", x.GetToken())
			}
		case ebnf.DynamicReference:
			{
				if consequent == -1 {
					for z := j - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == ebnf.Object {
							bin[j].object = bin[z].object
							break
						} else if bin[z].GetTokentype() == ebnf.Class {
							bin[j].class = bin[z].class
							break
						}
					}
				} else {
					for z := consequent - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == ebnf.DynamicReference && bin[z].GetToken() == x.GetToken() {
							bin[j].object = bin[z].object
							bin[j].class = bin[z].class
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
										ok = true
										break
									}
								}
							}
						}
					}
					if !ok {
						log.Fatal("Constant not found in KB! ", x.GetToken())
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

/*
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
*/
func (kb *KnowledgeBase) AddAttribute(c *KBClass, attrs ...*KBAttribute) {
	for i, _ := range attrs {
		attrs[i].Id = bson.NewObjectId()
		c.Attributes = append(c.Attributes, *attrs[i])
	}
	lib.LogFatal(c.Persist())
}

func (kb *KnowledgeBase) NewClass(c *KBClass) {
	for i, _ := range c.Attributes {
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
	for i, _ := range c.Attributes {
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
	for i, _ := range kb.Workspaces {
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
		n := KBAttributeObject{Id: bson.NewObjectId(), Attribute: x.Id, KbAttribute: x}
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

/*
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
*/

func (kb *KnowledgeBase) SaveValue(attr *KBAttributeObject, value any, source KBSource) *KBHistory {
	if attr != nil {
		h := KBHistory{Attribute: attr.Id, When: time.Now(), Value: value, Source: source}
		lib.LogFatal(h.Persist())
		attr.KbHistory = &h
		return &h
	} else {
		log.Fatal("Invalid Attribute of Object!")
		return nil
	}
}

func (kb *KnowledgeBase) FindClassByName(name string, mandatory bool) *KBClass {
	var ret KBClass
	err := ret.FindOne(bson.D{{"name", name}})
	if err != nil && mandatory {
		log.Fatal(err)
	}
	return kb.IdxClasses[ret.Id]
}

/*
	func (kb *KnowledgeBase) FindClassById(id bson.ObjectId) *KBClass {
		return kb.findClassBin(id, 0, len(kb.Classes)-1)
	}
*/
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
	log.Fatal(obj.Persist())
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

func (kb *KnowledgeBase) NewRule(rule string, priority byte, interval int) *KBRule {
	_, bin, err := kb.ParsingCommand(rule)
	if err != nil {
		log.Fatal(err)
	}
	r := KBRule{Rule: rule, Priority: priority, ExecutionInterval: interval}
	lib.LogFatal(r.Persist())
	kb.linkerRule(&r, bin)
	kb.Rules = append(kb.Rules, r)
	return &r
}

func (kb *KnowledgeBase) ReadBK() {
	log.Println("ReadBK")
	FindAllClasses("_id", &kb.Classes)
	for j, _ := range kb.Classes {
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

	FindAllObjects("name", &kb.Objects)
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
				err := h.FindLast(bson.D{{"attribute_id", x.Id}})
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

	FindAllWorkspaces("name", &kb.Workspaces)

	FindAllRules("_id", &kb.Rules)

	for i, _ := range kb.Rules {
		_, bin, err := kb.ParsingCommand(kb.Rules[i].Rule)
		if err != nil {
			log.Fatal(err)
		}
		kb.linkerRule(&kb.Rules[i], bin)
	}
}

func (kb *KnowledgeBase) ReadEBNF(file string) {
	ebnf := ebnf.EBNF{}
	kb.ebnf = &ebnf
	kb.ebnf.ReadToken(file)
}

func (kb *KnowledgeBase) PrintEBNF() {
	kb.ebnf.PrintEBNF()
}
