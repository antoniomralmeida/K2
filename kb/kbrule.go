package kb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/fuzzy"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *KBRule) String() string {
	j, err := json.MarshalIndent(*r, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}

func (obj *KBRule) Persist() error {
	return initializers.Persist(obj)

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
	initializers.Log("run..."+r.ID.Hex(), initializers.Info)

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
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			i++
			if r.bin[i].tokentype != ebnf.Class {
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			if r.bin[i].class == nil {
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token+" KB Class not found!", initializers.Error)
			}

			if len(r.bin[i].objects) == 0 {
				return initializers.Log("Warning in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token+" no object found!", initializers.Info)
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
					return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				i++
				if r.bin[i].tokentype != ebnf.Attribute {
					return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
				}

				if r.bin[i].class == nil {
					return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
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
						return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
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
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
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
					return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				if r.bin[i].attributeObjects == nil {
					return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
				}
				key := "{{" + r.bin[i].class.Name + "." + r.bin[i].token + "}}"
				txt = txt + " " + key + " "
				attrs[key] = r.bin[i].attributeObjects
				cart.AddItem(key, len(attrs[key])-1)

				i += 2
				if r.bin[i].literalbin == b_the {
					i += 2
				} else if r.bin[i].tokentype != ebnf.DynamicReference {
					return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
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

		case b_set:
			i += 2
			if r.bin[i].tokentype != ebnf.Attribute {
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			if r.bin[i].attributeObjects == nil {
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			attrs := r.bin[i].attributeObjects
			if r.bin[i+3].tokentype != ebnf.Literal && r.bin[i+4].tokentype != ebnf.Literal {
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
			}
			if r.bin[i+4].tokentype != ebnf.Constant && r.bin[i+5].tokentype != ebnf.Constant {
				return initializers.Log("Error in KB Rule "+r.ID.Hex()+" near "+r.bin[i].token, initializers.Error)
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
