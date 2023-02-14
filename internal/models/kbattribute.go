package models

import (
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KBAttribute struct {
	ID               primitive.ObjectID `bson:"id"`
	Name             string             `bson:"name" valid:"length(2|50),required"`
	AType            KBAttributeType    `bson:"atype" valid:"required"`
	KeepHistory      int                `bson:"keephistory" valid:range(0|5000)`             //Numero de historico a manter, 0- manter todos
	ValidityInterval int64              `bson:"validityinterval" valid:range(0|86400000000)` //validade do ultimo valor em microssegudos, 0- sempre
	SimulationID     KBSimulation       `bson:"simulation,omitempty" json:"-"`
	Simulation       string             `bson:"-" json:"simulation"`
	SourcesID        []KBSource         `bson:"sources" valid:"required"`
	Options          []string           `bson:"options,omitempty"`
	Sources          []string           `bson:"-" json:"sources"`
	antecedentRules  []*KBRule          `bson:"-"`
	consequentRules  []*KBRule          `bson:"-"`
}

func (obj *KBAttribute) Valitate() (bool, error) {
	return govalidator.ValidateStruct(obj)
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
