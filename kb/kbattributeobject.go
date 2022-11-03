package kb

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/montanaflynn/stats"
	"gonum.org/v1/gonum/stat"
	"gopkg.in/mgo.v2/bson"
)

func (ao *KBAttributeObject) Validity() bool {
	if ao.KbHistory != nil {
		if ao.KbAttribute.Deadline != 0 {
			diff := time.Now().Sub(time.Unix(0, ao.KbHistory.When))
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
		if KBDate == ao.KbAttribute.AType {
			i, _ := strconv.ParseInt(fmt.Sprintf("%v", ao.KbHistory.Value), 10, 64)
			return time.Unix(0, i)
		}
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
func (attr *KBAttributeObject) getFullName() string {
	return attr.KbObject.Name + "." + attr.KbAttribute.Name
}

func (attr *KBAttributeObject) String() string {
	j, err := json.MarshalIndent(*attr, "", "\t")
	lib.LogFatal(err)
	return string(j)
}

func (attr *KBAttributeObject) SetValue(value any, source KBSource, certainty float32) *KBHistory {
	if attr == nil {
		log.Println("Invalid attribute!")
		return nil
	}
	if !attr.KbAttribute.isSource(source) && source != Inference {
		log.Println("Invalid attribute source!")
		return nil
	}
	if reflect.TypeOf(value).String() == "string" {
		str := fmt.Sprintf("%v", value)
		switch attr.KbAttribute.AType {
		case KBNumber:
			value, _ = strconv.ParseFloat(str, 64)
		case KBDate:
			t, err := time.Parse("02/01/2006", str)
			if err == nil {
				value = t.UnixNano()
			} else {
				log.Println()
				return nil
			}
		}
	}
	h := KBHistory{Attribute: attr.Id, When: time.Now().UnixNano(), Value: value, Source: source, Trust: certainty}
	lib.LogFatal(h.Persist())
	attr.KbHistory = &h
	return &h
}

func (attr *KBAttributeObject) LinearRegression() error {
	log.Println("LinearRegression...")
	if attr.KbAttribute.AType == KBNumber {
		c := initializers.GetDb().C("KBHistory")
		pipe := c.Pipe([]bson.M{bson.M{"$match": bson.M{"attribute_id": attr.Id}}, bson.M{"$project": bson.M{"_id": 0, "value": 1, "when": 1, "trust": 1}}})
		resp := []bson.M{}
		iter := pipe.Iter()
		err := iter.All(&resp)
		lib.LogFatal(err)
		if len(resp) <= 2 {
			log.Println("cannot do linear regression with | C|<=2")
			return nil
		}
		X := make([]float64, len(resp))
		Y := make([]float64, len(resp))
		T := make([]float64, len(resp))
		for i := range resp {
			y, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[i]["value"]), 64)
			x, _ := strconv.ParseInt(fmt.Sprintf("%v", resp[i]["when"]), 10, 64)
			t, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[i]["trust"]), 32)

			X[i] = float64(x)
			Y[i] = y
			T[i] = t
		}
		trust := stat.Mean(T, nil) / 100.0
		alpha, beta := stat.LinearRegression(X, Y, nil, false)
		r2 := stat.RSquared(X, Y, nil, alpha, beta)
		xn := float64(time.Now().UnixNano())
		fx := alpha + xn*beta
		attr.SetValue(fx, KBSource(Simulation), float32(r2*trust*100.0))
	}
	return nil
}

func (attr *KBAttributeObject) MonteCarlo() error {
	log.Println("MonteCarlo...")
	fmt.Println("MonteCarlo...")
	if attr.KbAttribute.AType == KBNumber {
		c := initializers.GetDb().C("KBHistory")
		pipe := c.Pipe([]bson.M{bson.M{"$match": bson.M{"attribute_id": attr.Id}},
			bson.M{"$group": bson.M{"_id": "$attribute_id",
				"avg":    bson.M{"$avg": "$value"},
				"stdDev": bson.D{{"$stdDevPop", "$value"}},
				"trust":  bson.D{{"$avg", "$trust"}},
			}}})
		resp := []bson.M{}
		iter := pipe.Iter()
		err := iter.All(&resp)
		lib.LogFatal(err)
		avg, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[0]["avg"]), 64)
		stdDev, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[0]["stdDev"]), 64)
		trust, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[0]["trust"]), 32)
		r := stats.NormPpfRvs(avg, stdDev, 1)[0]
		a := stats.NormPpf(r, avg, stdDev) * (trust / 100.0)
		attr.SetValue(r, KBSource(Simulation), float32(a*100))
	}

	return nil
}

func (attr *KBAttributeObject) NormalDistribution() error {
	log.Println("NormalDistribution...")
	fmt.Println("NormalDistribution...")
	if attr.KbAttribute.AType == KBNumber {
		c := initializers.GetDb().C("KBHistory")
		pipe := c.Pipe([]bson.M{bson.M{"$match": bson.M{"attribute_id": attr.Id}},
			bson.M{"$group": bson.M{"_id": "$attribute_id",
				"avg":    bson.M{"$avg": "$value"},
				"stdDev": bson.D{{"$stdDevPop", "$value"}},
				"trust":  bson.D{{"$avg", "$trust"}},
			}}})
		resp := []bson.M{}
		iter := pipe.Iter()
		err := iter.All(&resp)
		lib.LogFatal(err)
		avg, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[0]["avg"]), 64)
		stdDev, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[0]["stdDev"]), 64)
		trust, _ := strconv.ParseFloat(fmt.Sprintf("%v", resp[0]["trust"]), 32)
		r := stats.NormPpfRvs(avg, stdDev, 1)[0]
		a := stats.NormPpf(r, avg, stdDev) * (trust / 100.0)
		attr.SetValue(r, KBSource(Simulation), float32(a*100))
	}

	return nil
}
