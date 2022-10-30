package kb

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	fmt.Println("run...", r.Rule)

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
