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
				//TODO: Acionar regras em backward
				//TODO: Create uma tarefa de simulação
				//TODO: As tarefas de busca de valor devem ter timite de tempo
				//TODO: Criar formulário web para receber valores de atributos de origem User (assincrono)
				//TODO: Levar em consideração a certeze na obteção de um valor PLC e User 100%, criar regra de envelhecimento da certeza
				//TODO: a certeza de um valor simulado deve analizer os quadrantes da curva normal do historico de valor
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
