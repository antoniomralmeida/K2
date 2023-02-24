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
	"github.com/asaskevich/govalidator"
	"github.com/kamva/mgm/v3"

	"github.com/antoniomralmeida/k2/internal/fuzzy"
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/pkg/cartesian"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KBRule struct {
	mgm.DefaultModel  `json:",inline" bson:",inline"`
	Name              string     `bson:"name" valid:"length(5|50),required"`
	Statement         string     `bson:"statement" valid:"required"`
	Priority          byte       `bson:"priority"` //0..100
	ExecutionInterval int        `bson:"interval"`
	Lastexecution     time.Time  `bson:"lastexecution"`
	consequent        int        `bson:"-"`
	inRun             bool       `bson:"-"`
	bkclasses         []*KBClass `bson:"-"`
	bin               []*BIN     `bson:"-"`
}

func (obj *KBRule) ValidateIndex() error {
	cur, err := mgm.Coll(obj).Indexes().List(mgm.Ctx())
	inits.Log(err, inits.Error)
	var result []bson.M
	err = cur.All(mgm.Ctx(), &result)
	if len(result) == 1 {
		inits.CreateUniqueIndex(mgm.Coll(obj), "name")
	}
	return err
}

func (obj *KBRule) validate() (bool, error) {
	return govalidator.ValidateStruct(obj)
}

func RuleFactory(name, statement string, priority byte, interval int) (*KBRule, error) {

	bin, err, detail := parsingRule(statement)
	if err != nil {
		inits.Log(fmt.Sprintf("%v %v", err, detail.String()), inits.Error)
		return nil, err
	}
	rule := KBRule{Name: name, Statement: statement, Priority: priority, ExecutionInterval: interval}

	ok, err := rule.validate()
	inits.Log(err, inits.Error)
	if !ok {
		return nil, err
	}
	err = linkerRule(&rule, bin)
	if err != nil {
		inits.Log(err, inits.Error)
		return nil, err
	}
	err = rule.Persist()
	if mongo.IsDuplicateKeyError(err) {
		inits.Log(err, inits.Error)
	} else {
		inits.Log(err, inits.Fatal)
	}
	if err == nil {
		return &rule, nil
	} else {
		return nil, err
	}
}

func (r *KBRule) String() string {
	j, err := json.MarshalIndent(*r, "", "\t")
	inits.Log(err, inits.Error)
	return string(j)
}

func (obj *KBRule) Persist() error {
	return inits.Persist(obj)
}
func FindRuleByName(name string) (ret *KBRule) {
	ret = nil
	cur := mgm.Coll(ret).FindOne(mgm.Ctx(), bson.D{{"name", name}})
	inits.Log(cur.Err(), inits.Error)
	if cur.Err() == nil {
		ret = new(KBRule)
		cur.Decode(ret)
	}
	return
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

func FindAllRules(sort string, rules *[]KBRule) error {
	collection := mgm.Coll(new(KBRule))
	cursor, err := collection.Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), rules)
	return err
}

func tokeningStatement(cmd string) []string {
	cmd = strings.Replace(cmd, "\r\n", "", -1)
	cmd = strings.Replace(cmd, "\\n", "", -1)
	cmd = strings.Replace(cmd, "\t", " ", -1)
	for strings.Contains(cmd, "  ") {
		cmd = strings.Replace(cmd, "  ", " ", -1)
	}
	inits.Log("Parsing Statement: "+cmd, inits.Info)
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
	return tokens
}

type compilingDetail struct {
	Token string
	Nexts map[string]*Token
}

func (cd *compilingDetail) String() string {
	expected := "( "
	for _, t := range cd.Nexts {
		expected += t.Token + " "
	}
	expected += ")"
	return fmt.Sprintf("near %v expected %v", cd.Token, expected)
}

func compilingStatement(tokens []string) ([]*BIN, error, compilingDetail) {
	pt := _ebnf.GetBase()
	jumps := []*Token{}
	nexts := make(map[string]*Token)
	bin := []*BIN{}
	for _, token := range tokens {
		var ok = false
		nexts = _ebnf.FindOptions(pt, &jumps, 0)
		for _, next := range nexts {
			if (next.GetToken() == token) ||
				(next.GetTokenType() == DynamicReference && len(token) == 1) ||
				((next.GetTokenType() == Object || next.GetTokenType() == Class || next.GetTokenType() == Rule ||
					next.GetTokenType() == Attribute || next.GetTokenType() == Constant ||
					next.GetTokenType() == Reference) && unicode.IsUpper(rune(token[0]))) ||
				(next.GetTokenType() == Text && (rune(token[0]) == '\'' || rune(token[0]) == '"') ||
					(next.GetTokenType() == Constant && lib.IsNumber(token))) {
				if next.GetTokenType() == Class {
					if FindClassByName(token, false) != nil {
						ok = true
					}
				} else if next.GetTokenType() == Object {
					if FindObjectByName(token) != nil {
						ok = true
					}
				} else if next.GetTokenType() == Rule {
					if FindRuleByName(token) != nil {
						ok = true
					}
				} else {
					ok = true
				}
				if ok {
					if pt.Rule_id != next.Rule_id {
						jumps = append(jumps, pt.Nexts...)
					}
					pt = next
					break
				}
			}
		}
		if !ok || len(nexts) == 0 {
			return nil, lib.CompilerError, compilingDetail{Token: token, Nexts: nexts}
		}
		code := BIN{TokenType: pt.GetTokenType(), Token: token}
		err := code.CheckLiteralBin()
		if err != nil {
			return nil, inits.Log(err, inits.Error), compilingDetail{}
		}
		bin = append(bin, &code)
	}
	nexts = _ebnf.FindOptions(pt, &jumps, 0)
	for _, y := range nexts {
		if y.GetToken() == "." && y.GetTokenType() == Control {
			return bin, nil, compilingDetail{}
		}
	}
	return nil, lib.CompilerError, compilingDetail{Token: "<end>"}
}

func parsingRule(cmd string) ([]*BIN, error, compilingDetail) {
	return compilingStatement(tokeningStatement(cmd))
}

func linkerRule(r *KBRule, bin []*BIN) error {
	// Find references of objects in KB
	inits.Log("Linking Production Rule: "+r.Name, inits.Info)
	pauseKB()

	dr := make(map[string]*KBClass)
	consequent := -1
	for j, x := range bin {
		switch x.LiteralBin {
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
				bin[j].workspace = FindWorkspaceByName(r.bin[j].Token)
			}
		case Object:
			if len(bin[j].objects) == 0 {
				obj := FindObjectByName(bin[j].Token)
				if obj != nil {
					bin[j].objects = append(bin[j].objects, obj)
				}
			}
		case Class:
			if bin[j].class == nil {
				c := FindClassByName(x.GetToken(), true)
				if c == nil {
					return lib.ClassNotFoundError
				}
				bin[j].class = c
				objs, err := FindAllObjects(bson.M{"class_id": c.ID}, "_id")
				inits.Log(err, inits.Error)
				for _, y := range objs {
					bin[j].objects = append(bin[j].objects, &y)
				}
			}
		case Attribute:
			ref := -1
			if bin[j+1].LiteralBin == B_of {
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
						obj := FindObjectByName(r.bin[j].Token)
						bin[j].objects = append(bin[j].objects, obj)
						bin[j].class = obj.ClassPtr
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
					objs, err := FindAllObjects(bson.M{"class_id": c.ID}, "_id")
					inits.Log(err, inits.Fatal)
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
						c = dr[bin[ref].Token]
						bin[ref].class = c
					}
					if c == nil {
						return inits.Log("Attribute class not found in KB! "+x.GetToken(), inits.Error)
					}
					bin[j].attribute = c.FindAttribute(x.GetToken())
					objs, err := FindAllObjects(bson.M{"class_id": c.ID}, "_id")
					inits.Log(err, inits.Fatal)
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
							dr[x.Token] = bin[j].class
							break
						}
					}
				} else {
					for z := consequent - 1; z >= 0; z-- {
						if bin[z].GetTokentype() == DynamicReference && bin[z].GetToken() == x.GetToken() {
							bin[j].objects = bin[z].objects
							bin[j].class = bin[z].class
							dr[x.Token] = bin[j].class
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
										bin[j].Token = "\"" + bin[j].Token + "\""
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
			bin[j].objects[z].Parsed = true
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
		if !_objects[i].Parsed {
			for j := range _rules {
				for k := range _rules[j].bkclasses {
					if _rules[j].bkclasses[k] == _objects[i].ClassPtr {
						bin, err, _ := parsingRule(_rules[j].Statement)
						if inits.Log(err, inits.Error) != nil {
							linkerRule(&_rules[j], bin)
						}
					}
				}
			}
			_objects[i].Parsed = true
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
	inits.Log("run..."+r.Name, inits.Info)

	attrs := make(map[string][]*KBAttributeObject)
	objs := make(map[string][]*KBObject)

	trueAntecedent := false
	expression := ""
	fuzzyexp := ""

oulter:
	//Program counter [pc] – It stores the counter which contains the address of the next instruction that is to be executed for the process.
	for pc := 0; pc < len(r.bin); {
		switch r.bin[pc].LiteralBin {
		case B_unconditionally:
			trueAntecedent = true
			pc++
		case B_then:
			if trueAntecedent {
				break oulter
			}
		case B_for:
			pc++
			if r.bin[pc].LiteralBin != B_any {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token, inits.Error)
			}
			pc++
			if r.bin[pc].TokenType != Class {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token, inits.Error)
			}
			if r.bin[pc].class == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token+" KB Class not found!", inits.Error)
			}

			if len(r.bin[pc].objects) == 0 {
				return inits.Log("Warning in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token+" no object found!", inits.Info)
			}

			if r.bin[pc+1].TokenType == DynamicReference {
				pc++
			}
		case B_if:

		inner:
			for {

				pc++
				for ; r.bin[pc].LiteralBin == B_open_par; pc++ {
					expression = expression + r.bin[pc].Token
					fuzzyexp = fuzzyexp + r.bin[pc].Token
				}
				if r.bin[pc].LiteralBin != B_the {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token, inits.Error)
				}
				pc++
				if r.bin[pc].TokenType != Attribute {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token, inits.Error)
				}

				if r.bin[pc].class == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token, inits.Error)
				}
				key := "{{" + r.bin[pc].class.Name + "." + r.bin[pc].Token + "}}"
				expression = expression + key
				fuzzyexp = fuzzyexp + key
				attrs[key] = r.bin[pc].attributeObjects
				objs[key] = r.bin[pc].objects

				pc++
				if r.bin[pc].LiteralBin == B_of {
					pc++
					if r.bin[pc].TokenType != DynamicReference && r.bin[pc].TokenType != Object {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].Token, inits.Error)
					}
					pc++
				}
				switch r.bin[pc].LiteralBin {
				case B_is:
					expression = expression + "=="
				case B_equal:
					expression = expression + "=="
				case B_different:
					expression = expression + "!="
				case B_less:
					expression = expression + "<"
					pc += 2
					if r.bin[pc].LiteralBin == B_or {
						expression = expression + "="
						pc += 2
					}
				case B_greater:
					expression = expression + ">"
					pc += 2
					if r.bin[pc].LiteralBin == B_or {
						expression = expression + "="
						pc += 2
					}
				}
				pc++
				if r.bin[pc].TokenType == Constant || r.bin[pc].TokenType == Text || r.bin[pc].TokenType == ListType {
					expression = expression + r.bin[pc].Token
				}
				pc++
				for ; r.bin[pc].LiteralBin == B_close_par; pc++ {
					expression = expression + r.bin[pc].Token
					fuzzyexp = fuzzyexp + r.bin[pc].Token
				}

				switch r.bin[pc].LiteralBin {
				case B_then:
					break inner
				case B_and:
					pc++
					expression = expression + " " + r.bin[pc].Token + " "
					fuzzyexp = fuzzyexp + " " + r.bin[pc].Token + " "
				case B_or:
					pc++
					fuzzyexp = fuzzyexp + " " + r.bin[pc].Token + " "
				}
			}
		default:
			pc++
		}
	}

	if !trueAntecedent {
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

type CreateType byte

const (
	NewClass CreateType = iota
	CopyClass
	InheritClass
	NewInstance
)

// TODO: ALTERAR COMAND SET
func RunDetailError(rule string, bin *BIN, pc int) error {
	return errors.New("Error in run rule " + rule + " near " + bin.Token + "[" + strconv.Itoa(pc) + "]!")
}

func (r *KBRule) RunConsequent(objs []*KBObject, trust float64) error {
	//Program counter [pc] – It stores the counter which contains the address of the next instruction that is to be executed for the process.

	for pc := r.consequent; pc < len(r.bin); pc++ {
		switch r.bin[pc].LiteralBin {
		case B_inform:
			attrs := make(map[string][]*KBAttributeObject)
			cart := cartesian.Cartesian{}
			pc += 5
			if r.bin[pc].TokenType != Text {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			txt := ""
			ok := true
			for {
				txt = txt + r.bin[pc].Token
				pc++
				if r.bin[pc].LiteralBin != B_the {
					break
				}
				if r.bin[pc].TokenType != Attribute {
					return RunDetailError(r.Name, r.bin[pc], pc)
				}
				if r.bin[pc].attributeObjects == nil {
					return RunDetailError(r.Name, r.bin[pc], pc)
				}
				key := "{{" + r.bin[pc].class.Name + "." + r.bin[pc].Token + "}}"
				txt = txt + " " + key + " "
				attrs[key] = r.bin[pc].attributeObjects
				cart.AddItem(key, len(attrs[key])-1)

				pc += 2
				if r.bin[pc].LiteralBin == B_the {
					pc += 2
				} else if r.bin[pc].TokenType != DynamicReference {
					return RunDetailError(r.Name, r.bin[pc], pc)
				} else {
					if !attrs[key][pc].InObjects(objs) {
						ok = false
					}
					pc++
				}
				if r.bin[pc].TokenType == Text {
					txt = txt + " " + r.bin[pc].Token
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
			if r.bin[pc].TokenType != Attribute {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			if r.bin[pc].attributeObjects == nil {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			attrs := r.bin[pc].attributeObjects
			pc += 4
			if r.bin[pc].TokenType != Constant && r.bin[pc+1].TokenType != Constant {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			if r.bin[pc+1].TokenType == Constant {
				pc++
			}
			v := r.bin[pc].Token
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
			var createType CreateType
			pc++
			if r.bin[pc].TokenType != Literal {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			switch r.bin[pc].LiteralBin {
			case B_a: //Class
				pc += 2
				switch r.bin[pc].LiteralBin {
				case B_by:
					pc += 2
					if r.bin[pc].class == nil {
						return RunDetailError(r.Name, r.bin[pc], pc)
					}
					createType = CopyClass
					baseClass = r.bin[pc].class
					pc++
				case B_whose:
					pc += 3
					if r.bin[pc].class == nil {
						return RunDetailError(r.Name, r.bin[pc], pc)
					}
					createType = InheritClass
					baseClass = r.bin[pc].class
					pc++
				}
			case B_an: //Instance
				pc += 4
				if r.bin[pc].class == nil {
					return RunDetailError(r.Name, r.bin[pc], pc)
				}
				createType = NewInstance
				baseClass = r.bin[pc].class
			default:
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			if r.bin[pc].LiteralBin == B_named {
				pc += 2
				name := lib.Identify(r.bin[pc].Token)
				switch createType {
				case NewClass:
					if _, err := KBClassFactory(name, "", ""); err != nil {
						return RunDetailError(r.Name, r.bin[pc], pc)
					}
				case CopyClass:
					if _, err := KBClassCopy(name, baseClass); err != nil {
						return RunDetailError(r.Name, r.bin[pc], pc)
					}
				case InheritClass:
					if _, err := KBClassFactoryParent(name, "", parentClass); err != nil {
						return RunDetailError(r.Name, r.bin[pc], pc)
					}
				case NewInstance:
					if _, err := ObjectFactoryByClass(name, baseClass); err != nil {
						return RunDetailError(r.Name, r.bin[pc], pc)
					}
				}
			}
			pc++
		case B_conclude:
			pc += 6
			if len(r.bin[pc].attributeObjects) != 1 {
				return RunDetailError(r.Name, r.bin[pc], pc)
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
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			obj := r.bin[pc].objects[0]
			pc += 2
			if r.bin[pc].workspace == nil {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			w := r.bin[pc].workspace
			w.AddObject(obj, 0, 0)
		case B_alter:
			pc++
			if r.bin[pc].class == nil {
				return RunDetailError(r.Name, r.bin[pc], pc)
			}
			//alterClass := r.bin[pc].class
			pc++
			for r.bin[pc].LiteralBin == B_add {
				pc++
				//attributeName := r.bin[pc].token
				options := []string{}
				pc += 2
				atype := r.bin[pc].Token
				if KBattributeTypeStr(atype) == KBList {
					pc++
					for r.bin[pc].LiteralBin != B_close_par {
						pc++
						options = append(options, r.bin[pc].Token)
						pc++
					}
				}

			}

		default:
			return RunDetailError(r.Name, r.bin[pc], pc)
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
