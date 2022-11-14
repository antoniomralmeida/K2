package kb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PaesslerAG/gval"
	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/fuzzy"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *KBRule) String() string {
	j, err := json.MarshalIndent(*r, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}

func (r *KBRule) Persist() error {
	ctx, collection := initializers.GetCollection("KBRule")
	if r.Id.IsZero() {
		r.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, r)
		return err
	} else {
		_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: r.Id}}, r)
		return err
	}
}

func FindAllRules(sort string) error {
	ctx, collection := initializers.GetCollection("KBRule")
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, &GKB.Rules)
	return err
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

func (kb *KnowledgeBased) ParsingCommand(cmd string) ([]*ebnf.Token, []*BIN, error) {
	cmd = strings.Replace(cmd, "\r\n", "", -1)
	cmd = strings.Replace(cmd, "\\n", "", -1)
	cmd = strings.Replace(cmd, "\t", " ", -1)
	for strings.Contains(cmd, "  ") {
		cmd = strings.Replace(cmd, "  ", " ", -1)
	}
	initializers.Log("Parsing Prodution Rule: "+cmd, initializers.Info)
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
			initializers.Log(", compilation successfully!", initializers.Info)
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
	initializers.Log("Linking Prodution Rule: "+r.Id.Hex(), initializers.Info)

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
				initializers.Log(FindAllObjects(bson.M{"class_id": c.Id}, "_id", &objs), initializers.Error)
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
					initializers.Log(FindAllObjects(bson.M{"class_id": c.Id}, "_id", &objs), initializers.Fatal)
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
						return initializers.Log("Attribute class not found in KB! "+x.GetToken(), initializers.Error)
					}
					bin[j].attribute = kb.FindAttribute(c, x.GetToken())
					objs := []KBObject{}
					initializers.Log(FindAllObjects(bson.M{"class_id": c.Id}, "_id", &objs), initializers.Fatal)
					for _, y := range objs {
						obj := kb.IdxObjects[y.Name]
						bin[j].objects = append(bin[j].objects, obj)
						atro := kb.FindAttributeObject(obj, x.GetToken())
						bin[j].attributeObjects = append(bin[j].attributeObjects, atro)
					}
					break
				}
			} else {
				return initializers.Log("Attribute not found in KB! "+x.GetToken(), initializers.Error)
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
						return initializers.Log("List option not found in KB! "+x.GetToken(), initializers.Error)
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

func (r *KBRule) Run() (e error) {

	type Value struct {
		value string
		trust float64
		atype KBAttributeType
	}
	GKB.mutex.Lock()
	if r.inRun { //avoid non-parallel execution of the same rule
		GKB.mutex.Unlock()
		return
	}
	r.inRun = true
	GKB.mutex.Unlock()
	initializers.Log("run..."+r.Id.Hex(), initializers.Info)

	attrs := make(map[string][]*KBAttributeObject)
	objs := make(map[string][]*KBObject)

	conditionally := false
	expression := ""
	fuzzyexp := ""
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
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			i++
			if r.bin[i].tokentype != ebnf.Class {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			if r.bin[i].class == nil {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token+" KB Class not found!", initializers.Error)
			}

			if len(r.bin[i].objects) == 0 {
				return initializers.Log("Warning in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token+" no object found!", initializers.Info)
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
					fuzzyexp = fuzzyexp + r.bin[i].token
				}
				if r.bin[i].literalbin != b_the {
					return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				i++
				if r.bin[i].tokentype != ebnf.Attribute {
					return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
				}

				if r.bin[i].class == nil {
					return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				key := "{{" + r.bin[i].class.Name + "." + r.bin[i].token + "}}"
				expression = expression + key
				fuzzyexp = fuzzyexp + key
				attrs[key] = r.bin[i].attributeObjects
				objs[key] = r.bin[i].objects

				i++
				if r.bin[i].literalbin == b_of {
					i++
					if r.bin[i].tokentype != ebnf.DynamicReference && r.bin[i].tokentype != ebnf.Object {
						return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
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
					fuzzyexp = fuzzyexp + r.bin[i].token
				}

				switch r.bin[i].literalbin {
				case b_then:
					break inner
				case b_and:
					i++
					expression = expression + " " + r.bin[i].token + " "
					fuzzyexp = fuzzyexp + " " + r.bin[i].token + " "
				case b_or:
					i++
					fuzzyexp = fuzzyexp + " " + r.bin[i].token + " "
				}
			}
		default:
			i++
		}
	}

	if !conditionally {
		cart := lib.Cartesian{}
		values := make(map[string][]Value)
		idx2 := []string{}
		for ix := range attrs {
			vls := []Value{}
			cart.AddItem(ix, len(attrs[ix])-1)
			for iy := range attrs[ix] {
				v, t, at := attrs[ix][iy].ValueString()
				vls = append(vls, Value{v, t, at})
			}
			values[ix] = vls
			idx2 = append(idx2, ix)
		}

		for {
			exp := expression
			fuzzy := fuzzy.FuzzyLogicalInference(fuzzyexp)
			found, idxs := cart.GetCombination()
			obs := []*KBObject{}
			ok := true
			for key := range attrs {
				if values[key][idxs[key]].value != "" {
					ok = false
					break
				}
				exp = strings.Replace(exp, key, string(values[key][idxs[key]].value), -1)
				trust := fmt.Sprint(values[key][idxs[key]].trust)
				fuzzy = strings.Replace(fuzzy, key, trust, -1)
				obs = append(obs, objs[key][idxs[key]])
			}
			if ok {
				result, err := gval.Evaluate(exp, nil)
				initializers.Log(err, initializers.Error)
				trust, err := gval.Evaluate(fuzzy, nil)
				initializers.Log(err, initializers.Error)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%v", trust), 64)
				if result == true {
					r.RunConsequent(obs, t)
				}
			}
			if !found {
				break
			}
		}
	} else {
		r.RunConsequent([]*KBObject{}, 100.0)
	}
	r.lastexecution = time.Now()
	GKB.mutex.Lock()
	r.inRun = false
	GKB.mutex.Unlock()
	return nil
}

func (r *KBRule) RunConsequent(objs []*KBObject, trust float64) error {
	for i := r.consequent; i < len(r.bin); {
		switch r.bin[i].literalbin {
		case b_inform:
			attrs := make(map[string][]*KBAttributeObject)
			cart := lib.Cartesian{}
			i += 5
			if r.bin[i].tokentype != ebnf.Text {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			txt := ""
			ok := true
			for {
				txt = txt + r.bin[i].token
				i++
				if r.bin[i].literalbin != b_the {
					break
				}
				if r.bin[i].tokentype != ebnf.Attribute {
					return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				if r.bin[i].attributeObjects == nil {
					return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				key := "{{" + r.bin[i].class.Name + "." + r.bin[i].token + "}}"
				txt = txt + " " + key + " "
				attrs[key] = r.bin[i].attributeObjects
				cart.AddItem(key, len(attrs[key])-1)

				i += 2
				if r.bin[i].literalbin == b_the {
					i += 2
				} else if r.bin[i].tokentype != ebnf.DynamicReference {
					return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
				} else {
					if !attrs[key][i].InObjects(objs) {
						ok = false
					}
					i++
				}
				if r.bin[i].tokentype == ebnf.Text {
					txt = txt + " " + r.bin[i].token
					i++
				}
			}
			if ok {
				txtout := txt
				found, idxs := cart.GetCombination()
				wks := make(map[primitive.ObjectID]*KBWorkspace)
				for key := range attrs {
					ao := attrs[key][idxs[key]]
					value, _, _ := ao.ValueString()
					txtout = strings.Replace(txtout, key, value, -1)
					ws := ao.KbObject.GetWorkspaces()
					for w := range ws {
						wks[ws[w].Id] = ws[w]
					}
				}
				for k := range wks {
					wks[k].Posts.Enqueue(txtout)
				}
				if !found {
					break
				}
			}

		case b_set:
			i += 2
			if r.bin[i].tokentype != ebnf.Attribute {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			if r.bin[i].attributeObjects == nil {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			attrs := r.bin[i].attributeObjects
			if r.bin[i+3].tokentype != ebnf.Literal && r.bin[i+4].tokentype != ebnf.Literal {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			if r.bin[i+4].tokentype != ebnf.Constant && r.bin[i+5].tokentype != ebnf.Constant {
				return initializers.Log("Error in KB Rule "+r.Id.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			var v string
			if r.bin[i+4].tokentype == ebnf.Constant {
				i += 4
				v = r.bin[i].token
			} else {
				i += 5
				v = r.bin[i].token
			}
			for _, a := range attrs {
				for _, o := range objs {
					if a.KbObject == o {
						a.SetValue(v, Inference, trust)
					}
				}
			}
			i++
		case b_halt:
			GKB.halt = true

		}
		//TODO: create
		//TODO: transfer
		//TODO: delete
		//TODO: insert
		//TODO: remove
		//TODO: change
		//TODO: move
		//TODO: rotate
		//TODO: show
		//TODO: hide
		//TODO: activate
		//TODO: deactivate
		//TODO: focus
		//TODO: invoke
		//TODO: conclude

	}
	return nil
}
