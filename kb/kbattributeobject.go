package kb

import (
	"log"
	"time"

	"github.com/antoniomralmeida/k2/lib"
)

func (ao *KBAttributeObject) Value() any {
	if ao.KbHistory != nil {
		if ao.KbAttribute.Deadline != 0 {
			diff := time.Now().Sub(ao.KbHistory.When)
			if diff.Milliseconds() > ao.KbAttribute.Deadline {
				ao.KbHistory = nil
				return nil
			}
		}
		return ao.KbHistory.Value
	} else {
		return nil
	}
}

func (attr *KBAttributeObject) SaveValue(kb *KnowledgeBase, value any, source KBSource) *KBHistory {
	if kb != nil {
		h := KBHistory{Attribute: attr.Id, When: time.Now(), Value: value, Source: source}
		lib.LogFatal(h.Persist())
		attr.KbHistory = &h
		return &h
	} else {
		log.Fatal("Invalid KnowledgeBase!")
		return nil
	}
}
