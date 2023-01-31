package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/kamva/mgm/v3"

	"github.com/antoniomralmeida/k2/cmd/k2/ebnf"
	"github.com/antoniomralmeida/k2/cmd/k2/fuzzy"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KBRule struct {
	mgm.DefaultModel  `json:",inline" bson:",inline"`
	Rule              string     `bson:"rule"`
	Priority          byte       `bson:"priority"` //0..100
	ExecutionInterval int        `bson:"interval"`
	lastexecution     time.Time  `bson:"-"`
	consequent        int        `bson:"-"`
	inRun             bool       `bson:"-"`
	bkclasses         []*KBClass `bson:"-"`
	bin               []*BIN     `bson:"-"`
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
	GKB.mutex.Lock()
	if r.inRun { //avoid non-parallel execution of the same rule
		GKB.mutex.Unlock()
		return
	}
	r.inRun = true
	GKB.mutex.Unlock()
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
		case models.B_unconditionally:
			conditionally = true
		case models.B_then:
			if !conditionally {
				break oulter
			}
		case models.B_for:
			pc++
			if r.bin[pc].literalbin != models.B_any {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			pc++
			if r.bin[pc].tokentype != ebnf.Class {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc].class == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
			}

			if len(r.bin[pc].objects) == 0 {
				return inits.Log("Warning in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" no object found!", inits.Info)
			}

			if r.bin[pc+1].tokentype == ebnf.DynamicReference {
				pc++
			}
		case models.B_if:

		inner:
			for {

				pc++
				for ; r.bin[pc].literalbin == models.B_open_par; pc++ {
					expression = expression + r.bin[pc].token
					fuzzyexp = fuzzyexp + r.bin[pc].token
				}
				if r.bin[pc].literalbin != models.B_the {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				pc++
				if r.bin[pc].tokentype != ebnf.Attribute {
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
				if r.bin[pc].literalbin == models.B_of {
					pc++
					if r.bin[pc].tokentype != ebnf.DynamicReference && r.bin[pc].tokentype != ebnf.Object {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
					}
					pc++
				}
				switch r.bin[pc].literalbin {
				case models.B_is:
					expression = expression + "=="
				case models.B_equal:
					expression = expression + "=="
				case models.B_different:
					expression = expression + "!="
				case models.B_less:
					expression = expression + "<"
					pc += 2
					if r.bin[pc].literalbin == models.B_or {
						expression = expression + "="
						pc += 2
					}
				case models.B_greater:
					expression = expression + ">"
					pc += 2
					if r.bin[pc].literalbin == models.B_or {
						expression = expression + "="
						pc += 2
					}
				}
				pc++
				if r.bin[pc].tokentype == ebnf.Constant || r.bin[pc].tokentype == ebnf.Text || r.bin[pc].tokentype == ebnf.ListType {
					expression = expression + r.bin[pc].token
				}
				pc++
				for ; r.bin[pc].literalbin == models.B_close_par; pc++ {
					expression = expression + r.bin[pc].token
					fuzzyexp = fuzzyexp + r.bin[pc].token
				}

				switch r.bin[pc].literalbin {
				case models.B_then:
					break inner
				case models.B_and:
					pc++
					expression = expression + " " + r.bin[pc].token + " "
					fuzzyexp = fuzzyexp + " " + r.bin[pc].token + " "
				case models.B_or:
					pc++
					fuzzyexp = fuzzyexp + " " + r.bin[pc].token + " "
				}
			}
		default:
			pc++
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
	r.lastexecution = time.Now()
	GKB.mutex.Lock()
	r.inRun = false
	GKB.mutex.Unlock()
	return nil
}

func (r *KBRule) RunConsequent(objs []*KBObject, trust float64) error {
	//Program counter [pc] – It stores the counter which contains the address of the next instruction that is to be executed for the process.

	for pc := r.consequent; pc < len(r.bin); pc++ {
		switch r.bin[pc].literalbin {
		case models.B_inform:
			attrs := make(map[string][]*KBAttributeObject)
			cart := lib.Cartesian{}
			pc += 5
			if r.bin[pc].tokentype != ebnf.Text {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			txt := ""
			ok := true
			for {
				txt = txt + r.bin[pc].token
				pc++
				if r.bin[pc].literalbin != models.B_the {
					break
				}
				if r.bin[pc].tokentype != ebnf.Attribute {
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
				if r.bin[pc].literalbin == models.B_the {
					pc += 2
				} else if r.bin[pc].tokentype != ebnf.DynamicReference {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				} else {
					if !attrs[key][pc].InObjects(objs) {
						ok = false
					}
					pc++
				}
				if r.bin[pc].tokentype == ebnf.Text {
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
					ws := ao.KbObject.GetWorkspaces()
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

		case models.B_set:
			pc += 2
			if r.bin[pc].tokentype != ebnf.Attribute {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc].attributeObjects == nil {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			attrs := r.bin[pc].attributeObjects
			if r.bin[pc+3].tokentype != ebnf.Literal && r.bin[pc+4].tokentype != ebnf.Literal {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc+4].tokentype != ebnf.Constant && r.bin[pc+5].tokentype != ebnf.Constant {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			var v string
			if r.bin[pc+4].tokentype == ebnf.Constant {
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
		case models.B_create:
			var baseClass *KBClass
			var parentClass *KBClass
			createClass := false

			pc++
			if r.bin[pc].tokentype != ebnf.Literal {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			switch r.bin[pc].literalbin {
			case models.B_a: //Class
				pc++
				createClass = true
				if r.bin[pc].class == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
				}
				pc++
				if r.bin[pc].tokentype != ebnf.Literal {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}
				switch r.bin[pc].literalbin {
				case models.B_by:
					pc += 2
					if r.bin[pc].class == nil {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
					}
					baseClass = r.bin[pc].class
					pc++
				case models.B_whose:
					pc += 3
					if r.bin[pc].class == nil {
						return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
					}
					parentClass = r.bin[pc].class
					pc++
				case models.B_named:
				default:
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
				}

			case models.B_an: //Instance
				pc += 4
				if r.bin[pc].class == nil {
					return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token+" KB Class not found!", inits.Error)
				}
				baseClass = r.bin[pc].class
			default:
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			if r.bin[pc].literalbin == models.B_named {
				pc += 2
				if createClass {
					className := r.bin[pc].GetToken()
					if baseClass != nil {
						GKB.CopyClass(className, baseClass)
					} else {
						GKB.NewSimpleClass(className, parentClass)
					}
				} else {
					objectName := r.bin[pc].GetToken()
					GKB.NewSimpleObject(objectName, baseClass)
				}
			}
		case models.B_conclude:
			pc += 6
			if len(r.bin[pc].attributeObjects) != 1 {
				return inits.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[pc].token, inits.Error)
			}
			attributeObject := r.bin[pc].attributeObjects[0]
			pc += 2
			attributeObject.SetValue(r.bin[pc].GetToken(), Inference, trust)
		case models.B_halt:
			GKB.Pause()
			models.NewAlert(inits.I18n_halt, "") //All users
		case models.B_transfer:
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

		//TODO: focus
		//TODO: invoke

	}
	return nil
}
