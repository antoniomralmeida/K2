package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type DataInput struct {
	Name    string          `json:"name"`
	Atype   KBAttributeType `json:"atype"`
	Options []string        `json:"options"`
}

func KBGetDataInput() []*DataInput {
	objs := []KBObject{}
	mgm.Coll(new(KBObject)).SimpleFind(&objs, bson.D{})

	ret := []*DataInput{}
	for i := range objs {
		for j := range objs[i].Attributes {
			a := &objs[i].Attributes[j]
			if a.KbAttribute.isSource(FromUser) && !a.Validity() {
				di := DataInput{Name: a.KbObject.Name + "." + a.KbAttribute.Name, Atype: a.KbAttribute.AType, Options: a.KbAttribute.Options}
				ret = append(ret, &di)
			}
		}
	}
	return ret
}
