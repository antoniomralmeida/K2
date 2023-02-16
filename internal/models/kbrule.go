package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PaesslerAG/gval"
	"github.com/kamva/mgm/v3"

	"github.com/antoniomralmeida/k2/internal/fuzzy"
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/pkg/cartesian"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KBRule struct {
	mgm.DefaultModel  `json:",inline" bson:",inline"`
	Rule              string     `bson:"rule"`
	Priority          byte       `bson:"priority"` //0..100
	ExecutionInterval int        `bson:"interval"`
	Lastexecution     time.Time  `bson:"lastexecution"`
	consequent        int        `bson:"-"`
	inRun             bool       `bson:"-"`
	bkclasses         []*KBClass `bson:"-"`
	bin               []*BIN     `bson:"-"`
}

func RuleFactory(rule string, priority byte, interval int) *KBRule {
	_, bin, err := parsingRule(rule)
	if inits.Log(err, inits.Info) != nil {
		return nil
	}
	r := KBRule{Rule: rule, Priority: priority, ExecutionInterval: interval}
	inits.Log(r.Persist(), inits.Fatal)
	linkerRule(&r, bin)
	return &r
}

func (r *KBRule) String() string {
	j, err := json.MarshalIndent(*r, "", "\t")
	inits.Log(err, inits.Error)
	return string(j)
}

func (obj *KBRule) Persist() error {
	return inits.Persist(obj)

}

func (obj *KBRule) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
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

func (r *KBRule) Run() (e error) {

	type Value struct {
		value string
		trust float64
		atype KBAttributeType
	}

	if r.inRun { //avoid non-parallel execution of the same rule
		return
	}
	r.inRun = true
	inits.Log("run..."+r.ID.Hex(), inits.Info)

	attrs := make(map[string][]*KBAttributeObject)
	objs := make(map[string][]*KBObject)

	conditionally := false
	expression := ""
	fuzzyexp := ""

oulter:
	//Program counter [pc] – It stores the counter which contains the address of the next instruction that is to be executed for the process.
	for pc := 0; pc < len(r.bin); {
		switch r.bin[pc].literalbin {
		case B_unconditionally:
			conditionally = true
		case B_then:
			if !conditionally {
				break oulter
			}
		case B_for:
			pc++
			if r.bin[pc].literalbin != B_any {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			pc++
			if r.bin[pc].tokentype != Class {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc].class == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
			}

			if len(r.bin[pc].objects) == 0 {
				return inits.Log("Warning in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" no object found!", inits.Info)
			}

			if r.bin[pc+1].tokentype == DynamicReference {
				pc++
			}
		case B_if:

		inner:
			for {

				pc++
				for ; r.bin[pc].literalbin == B_open_par; pc++ {
					expression = expression + r.bin[pc].token
					fuzzyexp = fuzzyexp + r.bin[pc].token
				}
				if r.bin[pc].literalbin != B_the {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				pc++
				if r.bin[pc].tokentype != Attribute {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}

				if r.bin[pc].class == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				key := "{{" + r.bin[pc].class.Name + "." + r.bin[pc].token + "}}"
				expression = expression + key
				fuzzyexp = fuzzyexp + key
				attrs[key] = r.bin[pc].attributeObjects
				objs[key] = r.bin[pc].objects

				pc++
				if r.bin[pc].literalbin == B_of {
					pc++
					if r.bin[pc].tokentype != DynamicReference && r.bin[pc].tokentype != Object {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
					}
					pc++
				}
				switch r.bin[pc].literalbin {
				case B_is:
					expression = expression + "=="
				case B_equal:
					expression = expression + "=="
				case B_different:
					expression = expression + "!="
				case B_less:
					expression = expression + "<"
					pc += 2
					if r.bin[pc].literalbin == B_or {
						expression = expression + "="
						pc += 2
					}
				case B_greater:
					expression = expression + ">"
					pc += 2
					if r.bin[pc].literalbin == B_or {
						expression = expression + "="
						pc += 2
					}
				}
				pc++
				if r.bin[pc].tokentype == Constant || r.bin[pc].tokentype == Text || r.bin[pc].tokentype == ListType {
					expression = expression + r.bin[pc].token
				}
				pc++
				for ; r.bin[pc].literalbin == B_close_par; pc++ {
					expression = expression + r.bin[pc].token
					fuzzyexp = fuzzyexp + r.bin[pc].token
				}

				switch r.bin[pc].literalbin {
				case B_then:
					break inner
				case B_and:
					pc++
					expression = expression + " " + r.bin[pc].token + " "
					fuzzyexp = fuzzyexp + " " + r.bin[pc].token + " "
				case B_or:
					pc++
					fuzzyexp = fuzzyexp + " " + r.bin[pc].token + " "
				}
			}
		default:
			pc++
		}
	}

	if !conditionally {
		cart := cartesian.Cartesian{}
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
				inits.Log(err, inits.Error)
				trust, err := gval.Evaluate(fuzzy, nil)
				inits.Log(err, inits.Error)
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
	r.Lastexecution = time.Now()
	r.Persist()
	r.inRun = false
	return nil
}

func (r *KBRule) RunConsequent(objs []*KBObject, trust float64) error {
	//Program counter [pc] – It stores the counter which contains the address of the next instruction that is to be executed for the process.

	for pc := r.consequent; pc < len(r.bin); pc++ {
		switch r.bin[pc].literalbin {
		case B_inform:
			attrs := make(map[string][]*KBAttributeObject)
			cart := cartesian.Cartesian{}
			pc += 5
			if r.bin[pc].tokentype != Text {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			txt := ""
			ok := true
			for {
				txt = txt + r.bin[pc].token
				pc++
				if r.bin[pc].literalbin != B_the {
					break
				}
				if r.bin[pc].tokentype != Attribute {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				if r.bin[pc].attributeObjects == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				key := "{{" + r.bin[pc].class.Name + "." + r.bin[pc].token + "}}"
				txt = txt + " " + key + " "
				attrs[key] = r.bin[pc].attributeObjects
				cart.AddItem(key, len(attrs[key])-1)

				pc += 2
				if r.bin[pc].literalbin == B_the {
					pc += 2
				} else if r.bin[pc].tokentype != DynamicReference {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				} else {
					if !attrs[key][pc].InObjects(objs) {
						ok = false
					}
					pc++
				}
				if r.bin[pc].tokentype == Text {
					txt = txt + " " + r.bin[pc].token
					pc++
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
					ws := KBGetWorkspacesFromObject(ao.KbObject)
					for w := range ws {
						wks[ws[w].ID] = ws[w]
					}
				}
				for k := range wks {
					wks[k].Posts.Enqueue(txtout)
				}
				if !found {
					break
				}
			}

		case B_set:
			pc += 2
			if r.bin[pc].tokentype != Attribute {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc].attributeObjects == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			attrs := r.bin[pc].attributeObjects
			if r.bin[pc+3].tokentype != Literal && r.bin[pc+4].tokentype != Literal {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc+4].tokentype != Constant && r.bin[pc+5].tokentype != Constant {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			var v string
			if r.bin[pc+4].tokentype == Constant {
				pc += 4
				v = r.bin[pc].token
			} else {
				pc += 5
				v = r.bin[pc].token
			}
			for _, a := range attrs {
				for _, o := range objs {
					if a.KbObject == o {
						a.SetValue(v, Inference, trust)
					}
				}
			}
		case B_create:
			var baseClass *KBClass
			var parentClass *KBClass
			createClass := false

			pc++
			if r.bin[pc].tokentype != Literal {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			switch r.bin[pc].literalbin {
			case B_a: //Class
				pc++
				createClass = true
				if r.bin[pc].class == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
				}
				pc++
				if r.bin[pc].tokentype != Literal {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				switch r.bin[pc].literalbin {
				case B_by:
					pc += 2
					if r.bin[pc].class == nil {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
					}
					baseClass = r.bin[pc].class
					pc++
				case B_whose:
					pc += 3
					if r.bin[pc].class == nil {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
					}
					parentClass = r.bin[pc].class
					pc++
				case B_named:
				default:
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}

			case B_an: //Instance
				pc += 4
				if r.bin[pc].class == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
				}
				baseClass = r.bin[pc].class
			default:
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc].literalbin == B_named {
				pc += 2
				if createClass {
					className := r.bin[pc].GetToken()
					if baseClass != nil {
						KBClassCopy(className, baseClass)
					} else {
						KBClassFactoryParent(className, "", parentClass)
					}
				} else {
					objectName := r.bin[pc].GetToken()
					ObjectFactoryByClass(objectName, baseClass)
				}
			}
		case B_conclude:
			pc += 6
			if len(r.bin[pc].attributeObjects) != 1 {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			attributeObject := r.bin[pc].attributeObjects[0]
			pc += 2
			attributeObject.SetValue(r.bin[pc].GetToken(), Inference, trust)
		case B_halt:
			pauseKB()
			NewAlert(inits.I18n_halt, "") //All users
		case B_transfer:
			pc++
			if len(r.bin[pc].objects) == 0 {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
			}
			obj := r.bin[pc].objects[0]
			pc += 2
			if r.bin[pc].workspace == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
			}
			w := r.bin[pc].workspace
			w.AddObject(obj, 0, 0)
		case B_alter:
			pc++
			if r.bin[pc].class == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
			}
			//alterClass := r.bin[pc].class
			pc++
			for r.bin[pc].literalbin == B_add {
				pc++
				//attributeName := r.bin[pc].token
				options := []string{}
				pc += 2
				atype := r.bin[pc].token
				if KBattributeTypeStr(atype) == KBList {
					pc++
					for r.bin[pc].literalbin != B_close_par {
						pc++
						options = append(options, r.bin[pc].token)
						pc++
					}
				}

			}

		default:
			return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
		}

		//TODO: delete
		//TODO: insert
		//TODO: remove
		//TODO: change
		//TODO: move
		//TODO: rotate
		//TODO: show
		//TODO: hide
		//TODO: alter
		//TODO: focus
		//TODO: invoke

	}
	return nil
}

func FindAllRules(sort string) error {
	collection := mgm.Coll(new(KBRule))
	cursor, err := collection.Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), _rules)
	return err
}

func parsingRule(cmd string) ([]*Token, []*BIN, error) {
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
		case c == '春' || c == ' ' || _ebnf.FindSymbols(string(c), true) != -1:
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
	var pt = _ebnf.GetBase()
	var stack []*Token
	var opts []*Token
	var bin []*BIN
	for _, x := range tokens {
		var ok = false
		opts = _ebnf.FindOptions(pt, &stack, 0)
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
			str := "Compiler error in " + x + " when the expected was: "
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
	opts = _ebnf.FindOptions(pt, &stack, 0)
	str := "Incomplete sentence when the expected was: "
	for _, y := range opts {
		str = str + "... " + y.GetToken()
	}
	return opts, nil, errors.New(str)
}

func linkerRule(r *KBRule, bin []*BIN) error {
	// Find references of objects in KB
	inits.Log("Linking Prodution Rule: "+r.ID.Hex(), inits.Info)
	pauseKB()

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
				bin[j].workspace = FindWorkspaceByName(r.bin[j].token)
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
					bin[j].attribute = bin[ref].class.FindAttribute(x.GetToken())
					if len(bin[j].objects) > 0 {
						atro := FindAttributeObject(bin[ref].objects[0], x.GetToken())
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
					bin[j].attribute = c.FindAttribute(x.GetToken())
					objs := []KBObject{}
					inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Fatal)
					for _, y := range objs {
						obj := &y
						bin[j].objects = append(bin[j].objects, obj)
						atro := FindAttributeObject(obj, x.GetToken())
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
					bin[j].attribute = c.FindAttribute(x.GetToken())
					objs := []KBObject{}
					inits.Log(FindAllObjects(bson.M{"class_id": c.ID}, "_id", &objs), inits.Fatal)
					for _, y := range objs {
						obj := &y
						bin[j].objects = append(bin[j].objects, obj)
						atro := FindAttributeObject(obj, x.GetToken())
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
	resumeKB()
	_rules = append(_rules, *r)
	return nil
}

func RefreshRules() error {
	inits.Log("RefreshRules...", inits.Info)
	for i := range _objects {
		if !_objects[i].parsed {
			for j := range _rules {
				for k := range _rules[j].bkclasses {
					if _rules[j].bkclasses[k] == _objects[i].Bkclass {
						_, bin, err := parsingRule(_rules[j].Rule)
						if inits.Log(err, inits.Error) != nil {
							linkerRule(&_rules[j], bin)
						}
					}
				}
			}
			_objects[i].parsed = true
		}
	}
	return nil
}

func runStackRules() error {
	inits.Log("RunStackRules...", inits.Info)
	for i := range _rules {
		if _rules[i].ExecutionInterval != 0 && time.Now().After(_rules[i].Lastexecution.Add(time.Duration(_rules[i].ExecutionInterval)*time.Millisecond)) {
			stack := KBStack{RuleID: _rules[i].ID}
			stack.Persist()
		}
	}

	toRun := RunFromStack()
	for _, r := range toRun {
		r.Run()
	}

	return nil
}
