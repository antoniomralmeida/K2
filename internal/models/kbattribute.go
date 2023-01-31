package models

import "github.com/kamva/mgm/v3"

type KBAttribute struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string          `bson:"name"`
	AType            KBAttributeType `bson:"atype"`
	KeepHistory      int             `bson:"keephistory"`      //Numero de historico a manter, 0- sempre
	ValidityInterval int64           `bson:"validityinterval"` //validade do ultimo valor em microssegudos, 0- sempre
	SimulationID     KBSimulation    `bson:"simulation,omitempty" json:"-"`
	Simulation       string          `bson:"-" json:"simulation"`
	SourcesID        []KBSource      `bson:"sources"`
	Options          []string        `bson:"options,omitempty"`
	Sources          []string        `bson:"-" json:"sources"`
	antecedentRules  []*KBRule       `bson:"-"`
	consequentRules  []*KBRule       `bson:"-"`
}

func (a *KBAttribute) addAntecedentRules(r *KBRule) {
	found := false
	for i := range a.antecedentRules {
		if a.antecedentRules[i] == r {
			found = true
			break
		}
	}
	if !found {
		a.antecedentRules = append(a.antecedentRules, r)
	}
}

func (a *KBAttribute) addConsequentRules(r *KBRule) {
	found := false
	for i := range a.consequentRules {
		if a.consequentRules[i] == r {
			found = true
			break
		}
	}
	if !found {
		a.consequentRules = append(a.consequentRules, r)
	}
}

func (a *KBAttribute) isSource(s KBSource) bool {
	for i := range a.SourcesID {
		if a.SourcesID[i] == s {
			return true
		}
	}
	return false
}
