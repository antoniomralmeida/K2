package kb

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/antoniomralmeida/k2/lib"
)

func (ao *KBAttributeObject) Validity() bool {
	if ao.KbHistory != nil {
		if ao.KbAttribute.Deadline != 0 {
			diff := time.Now().Sub(ao.KbHistory.When)
			if diff.Milliseconds() > ao.KbAttribute.Deadline {
				ao.KbHistory = nil
				return false
			}
		}
		return true
	}
	return false
}

func (ao *KBAttributeObject) Value() any {
	if ao.Validity() {
		return ao.KbHistory.Value
	} else {
		return nil
	}

	//TODO: Acionar regras em backward chaning
	//TODO: Criar uma tarefa de simulação
	//TODO: As tarefas de busca de valor devem ter limite de tempo
	//TODO: Criar formulário web para receber valores de atributos de origem User (assincrono)
	//TODO: Levar em consideração a certeza na obteção de um valor PLC e User 100%
	//TODO: Criar regra de envelhecimento da certeza, com base na disperção e na validade do dado
	//TODO: a certeza de um valor simulado deve analizar os quadrantes da curva normal do historico de valor
	//TODO: a certeza por inferencia deve usar logica fuzzi

}

func (attr *KBAttributeObject) SetValue(kb *KnowledgeBase, value any, source KBSource, certainty float32) *KBHistory {
	if kb != nil {

		if reflect.TypeOf(value).String() == "string" {
			str := fmt.Sprintf("%v", value)
			switch attr.KbAttribute.AType {
			case KBNumber:
				value, _ = strconv.ParseFloat(str, 64)
			case KBDate:
				value, _ = time.Parse("02/01/2006", str)
			}
		}
		h := KBHistory{Attribute: attr.Id, When: time.Now(), Value: value, Source: source, Certainty: certainty}
		lib.LogFatal(h.Persist())
		attr.KbHistory = &h
		return &h
	} else {
		log.Fatal("Invalid KnowledgeBase!")
		return nil
	}
}
