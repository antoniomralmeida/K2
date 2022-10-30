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
				//TODO: Acionar regras em backward chaning
				//TODO: Criar uma tarefa de simulação
				//TODO: As tarefas de busca de valor devem ter limite de tempo
				//TODO: Criar formulário web para receber valores de atributos de origem User (assincrono)
				//TODO: Levar em consideração a certeza na obteção de um valor PLC e User 100%
				//TODO: Criar regra de envelhecimento da certeza, com base na disperção e na validade do dado
				//TODO: a certeza de um valor simulado deve analizar os quadrantes da curva normal do historico de valor
				//TODO: a certeza por inferencia deve usar logica fuzzi

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
