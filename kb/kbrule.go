package kb

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/antoniomralmeida/k2/db"
	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/apaxa-go/eval"
	"gopkg.in/mgo.v2/bson"
)

func (r *KBRule) Run() {

	type ctrl struct {
		i   int
		max int
	}
	log.Println("run...", r.Id)

	attrs := make(map[string][]*KBAttributeObject)
	objs := make(map[string][]*KBObject)

	conditionally := false
	expression := ""
oulter:
	for i := 0; i < len(r.bin); {
		switch r.bin[i].literalbin {
		case b_unconditionally:
			conditionally = true
		case b_then:
			if !conditionally {
				break oulter
			}
		case b_for:
			i++
			if r.bin[i].literalbin != b_any {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
			}
			i++
			if r.bin[i].tokentype != ebnf.Class {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
			}
			if r.bin[i].class == nil {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token, " KB Class not found!")
			}

			if len(r.bin[i].objects) == 0 {
				log.Println("Warning in KB Rule ", r.Id, " near ", r.bin[i].token, " no object found!")
				break
			}

			if r.bin[i+1].tokentype == ebnf.DynamicReference {
				i++
			}
		case b_if:

		inner:
			for {

				i++
				for ; r.bin[i].literalbin == b_open_par; i++ {
					expression = expression + r.bin[i].token
				}
				if r.bin[i].literalbin != b_the {
					log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
				}
				i++
				if r.bin[i].tokentype != ebnf.Attribute {
					log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
				}

				if r.bin[i].class == nil {
					log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
				}
				key := "{{" + r.bin[i].class.Name + "." + r.bin[i].token + "}}"
				expression = expression + key
				attrs[key] = r.bin[i].attributeObjects
				objs[key] = r.bin[i].objects

				i++
				if r.bin[i].literalbin == b_of {
					i++
					if r.bin[i].tokentype != ebnf.DynamicReference && r.bin[i].tokentype != ebnf.Object {
						log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
					}
					i++
				}
				switch r.bin[i].literalbin {
				case b_is:
					expression = expression + "=="
				case b_equal:
					expression = expression + "=="
				case b_different:
					expression = expression + "!="
				case b_less:
					expression = expression + "<"
					i += 2
					if r.bin[i].literalbin == b_or {
						expression = expression + "="
						i += 2
					}
				case b_greater:
					expression = expression + ">"
					i += 2
					if r.bin[i].literalbin == b_or {
						expression = expression + "="
						i += 2
					}
				}
				i++
				if r.bin[i].tokentype == ebnf.Constant || r.bin[i].tokentype == ebnf.Text || r.bin[i].tokentype == ebnf.ListType {
					expression = expression + r.bin[i].token
				}
				i++
				for ; r.bin[i].literalbin == b_close_par; i++ {
					expression = expression + r.bin[i].token
				}

				switch r.bin[i].literalbin {
				case b_then:
					break inner
				case b_and:
					i++
					expression = expression + " " + r.bin[i].token + " "
				case b_or:
					i++
					expression = expression + " " + r.bin[i].token + " "
				}
			}
		default:
			i++
		}
	}

	if !conditionally {
		values := make(map[string][]any)
		idx := make(map[string]ctrl)
		idx2 := []string{}
		for ix := range attrs {
			vls := []any{}
			idx[ix] = ctrl{0, len(attrs[ix]) - 1}
			for iy := range attrs[ix] {
				value := attrs[ix][iy].Value()
				vls = append(vls, value)
			}
			values[ix] = vls
			idx2 = append(idx2, ix)
		}
		iz := 0
	i00:
		for {
			exp := expression
			obs := []*KBObject{}
			ok := true
			for key := range attrs {
				if values[key][idx[key].i] == nil {
					ok = false
					break
				}
				value := fmt.Sprint(values[key][idx[key].i])
				exp = strings.Replace(exp, key, value, -1)
				obs = append(obs, objs[key][idx[key].i])
			}
			if ok {
				exp = "bool(" + exp + ")"
				fmt.Println(exp)
				expr, err := eval.ParseString(exp, "")
				lib.LogFatal(err)
				result, err := expr.EvalToInterface(nil)
				lib.LogFatal(err)
				if result == true {
					r.RunConsequent(obs)
				}
			}
		i01:
			for {
				ix := idx2[iz]
				if idx[ix].i < idx[ix].max {
					idx[ix] = ctrl{idx[ix].i + 1, idx[ix].max}
					break i01
				} else {
					if iz >= len(idx2)-1 {
						break i00
					}
					iz++
				}
			}
		}
	} else {
		r.RunConsequent([]*KBObject{})
	}

	r.lastexecution = time.Now()
}

func (r *KBRule) RunConsequent(objs []*KBObject) {
	for i := r.consequent; i < len(r.bin); {
		switch r.bin[i].literalbin {
		case b_inform:
			i += 4
			if r.bin[i].tokentype != ebnf.Text {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
			}
		case b_set:
			//TODO: Acionar regras em forward chaining
		}
	}
}

func (r *KBRule) Persist() error {
	collection := db.GetDb().C("KBRule")
	if r.Id == "" {
		r.Id = bson.NewObjectId()
		return collection.Insert(r)
	} else {
		return collection.UpdateId(r.Id, r)
	}
}

func FindAllRules(sort string, rs *[]KBRule) error {
	collection := db.GetDb().C("KBRule")
	return collection.Find(bson.M{}).Sort(sort).All(rs)
}

func (r *KBRule) addClass(c *KBClass) {
	found := false
	for i := range r.bkclasses {
		if r.bkclasses[i] == c {
			found = true
			break
		}
	}
	if !found {
		r.bkclasses = append(r.bkclasses, c)
	}
}

func (r *KBRule) GetBins() []*BIN {
	return r.bin
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

	dr := make(map[string]*KBClass)
	consequent := -1
	for j, x := range bin {
		switch x.literalbin {
		case b_initially:
			kb.mutex.Lock()
			kb.stack = append(kb.stack, r)
			kb.mutex.Unlock()
		case b_then:
			consequent = j
			r.consequent = j + 1
		}
		switch x.GetTokentype() {
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
				lib.LogFatal(FindAllObjects(bson.M{"class_id": c.Id}, "_id", &objs))
				for _, y := range objs {
					bin[j].objects = append(bin[j].objects, kb.IdxObjects[y.Name])
				}
			}
		case ebnf.Attribute:
			ref := -1
			if bin[j+1].literalbin == b_of {
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
					lib.LogFatal(FindAllObjects(bson.M{"class_id": c.Id}, "_id", &objs))
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
						log.Fatal("Attribute class not found in KB! ", x.GetToken())
					}
					bin[j].attribute = kb.FindAttribute(c, x.GetToken())
					objs := []KBObject{}
					lib.LogFatal(FindAllObjects(bson.M{"class_id": c.Id}, "_id", &objs))
					for _, y := range objs {
						obj := kb.IdxObjects[y.Name]
						bin[j].objects = append(bin[j].objects, obj)
						atro := kb.FindAttributeObject(obj, x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				} else {

				}
			} else {
				log.Fatal("Attribute not found in KB! ", x.GetToken())
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
						log.Fatal("Constant not found in KB! ", x.GetToken())
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
}
