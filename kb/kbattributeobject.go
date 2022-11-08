package kb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/montanaflynn/stats"
	p "github.com/rafaeljesus/parallel-fn"
	"gonum.org/v1/gonum/stat"
	"gopkg.in/mgo.v2/bson"
)

func (ao *KBAttributeObject) Validity() bool {
	if ao.KbHistory != nil {
		if ao.KbAttribute.ValidityInterval != 0 {
			diff := time.Now().Sub(time.Unix(0, ao.KbHistory.When))
			if diff.Milliseconds() > ao.KbAttribute.ValidityInterval {
				ao.KbHistory = nil
				return false
			}
		}
		return true
	}
	return false
}

func (ao *KBAttributeObject) Value() (any, float64) {
	if ao.Validity() {
		if KBDate == ao.KbAttribute.AType {
			i, _ := strconv.ParseInt(fmt.Sprintf("%v", ao.KbHistory.Value), 10, 64)
			return time.Unix(0, i), ao.KbHistory.Trust
		}
		return ao.KbHistory.Value, ao.KbHistory.Trust
	} else {
		timeout := time.After(1 * time.Second) // real-time search
		fn1 := func() error {
			for _, r := range ao.KbAttribute.consequentRules { //backward chaining
				r.Run()
				if ao.KbHistory != nil { //when find a value (stop)
					return nil
				}
			}
			return nil
		}
		fn2 := func() error {
			if ao.KbAttribute.isSource(Simulation) {
				switch ao.KbAttribute.SimulationID {
				case MonteCarlo:
					ao.MonteCarlo()
				case LinearRegression:
					ao.LinearRegression()
				case NormalDistribution:
					ao.NormalDistribution()
				}
			}
			return nil
		}
		fn3 := func() error {
			if ao.KbAttribute.isSource(IOT) {
				ao.IOTParsing()
			}
			return nil
		}

		//TODO: testar a execução paralela
		for {
			select {
			case e := <-p.Run(fn1, fn2, fn3):
				log.Println(e)
				return nil, 0
			case <-timeout:
				if ao.KbHistory != nil {
					return ao.KbHistory.Value, ao.KbHistory.Trust
				}
			}
		}
	}

}
func (attr *KBAttributeObject) getFullName() string {
	return attr.KbObject.Name + "." + attr.KbAttribute.Name
}

func (attr *KBAttributeObject) String() string {
	j, err := json.MarshalIndent(*attr, "", "\t")
	lib.LogFatal(err)
	return string(j)
}

func (attr *KBAttributeObject) SetValue(value any, source KBSource, trust float64) *KBHistory {
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
	h := KBHistory{Attribute: attr.Id, When: time.Now().UnixNano(), Value: value, Source: source, Trust: trust}
	lib.LogFatal(h.Persist())
	attr.KbHistory = &h
	attr.Kb.stack = append(attr.Kb.stack, attr.KbAttribute.antecedentRules...) //  forward chaining
	if attr.KbAttribute.KeepHistory != 0 {
		go h.ClearingHistory(attr.KbAttribute.KeepHistory)
	}
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
		attr.SetValue(fx, KBSource(Simulation), r2*trust*100.0)
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
		attr.SetValue(r, KBSource(Simulation), a*100)
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
		attr.SetValue(r, KBSource(Simulation), a*100)
	}

	return nil
}

func (attr *KBAttributeObject) IOTParsing() error {
	log.Println("IOTParsing...")
	if !attr.Validity() {
		api := os.Getenv("IOTMIDWARE")
		if attr.KbAttribute.isSource(KBSource(User)) && api != "" {
			iotapi := api + "?" + attr.getFullName()
			api := fiber.AcquireAgent()
			req := api.Request()
			req.Header.SetMethod("post")
			req.SetRequestURI(iotapi)
			if err := api.Parse(); err != nil {
				log.Println(err)
			} else {
				_, body, errs := api.Bytes()
				if errs != nil {
					attr.SetValue(string(body), IOT, 100.0)
				}
			}
		}
	}
	return nil
}
