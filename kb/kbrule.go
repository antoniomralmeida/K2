package kb

import (
	"fmt"
	"log"
	"time"

	"github.com/antoniomralmeida/k2/db"
	"github.com/antoniomralmeida/k2/ebnf"
	"gopkg.in/mgo.v2/bson"
)

func (r *KBRule) Run() {

	log.Println("run...", r.Id)

	dr := make(map[string]int)
	conditionally := false
	conclude := false
	for i := 0; i < len(r.bin); i++ {
		switch r.bin[i].typebin {
		case b_unconditionally:
			if conclude {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
			}
			conditionally = true
		case b_then:
			if !conditionally {
				break
			}
			conclude = true
		case b_for:
			if conclude {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
			}
			i++
			if r.bin[i].typebin != b_any {
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
				dr[r.bin[i].token] = i
			}
		case b_if:
			expression := ""
			if conclude {
				log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
			}
			for {

				i++
				for ; r.bin[i].typebin == b_open_par; i++ {
					expression = expression + r.bin[i].token
				}
				if r.bin[i].typebin != b_the {
					log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
				}
				i++
				if r.bin[i].tokentype != ebnf.Attribute {
					log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
				}
				expression = expression + "{{" + r.bin[i].token + "}}"
				i++
				if r.bin[i].typebin == b_of {
					i++
					if r.bin[i].tokentype != ebnf.DynamicReference && r.bin[i].tokentype != ebnf.Object {
						log.Fatal("Error in KB Rule ", r.Id, " near ", r.bin[i].token)
					}
					i++
				}
				switch r.bin[i].typebin {
				case b_is:
					expression = expression + "="
				case b_equal:
					expression = expression + "="
				case b_different:
					expression = expression + "!="
				case b_less:
					expression = expression + "<"
					i += 2
					if r.bin[i].typebin == b_or {
						expression = expression + "="
						i += 2
					}
				case b_greater:
					expression = expression + ">"
					i += 2
					if r.bin[i].typebin == b_or {
						expression = expression + "="
						i += 2
					}
				}
				i++
				if r.bin[i].tokentype == ebnf.Constant || r.bin[i].tokentype == ebnf.Text || r.bin[i].tokentype == ebnf.ListType {
					expression = expression + r.bin[i].token
				}
				i++
				for ; r.bin[i].typebin == b_close_par; i++ {
					expression = expression + r.bin[i].token
				}
				i++
				switch r.bin[i].typebin {
				case b_then:
					break
				case b_and:
					expression = expression + " " + r.bin[i].token + " "
				case b_or:
					expression = expression + " " + r.bin[i].token + " "
				}
			}
			fmt.Println(expression)

		}
	}
	r.lastexecution = time.Now()
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
