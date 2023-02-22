package models

import (
	"math/rand"
	"testing"
	"time"

	"github.com/kamva/mgm/v3"

	"github.com/antoniomralmeida/k2/internal/inits"

	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	gotenv.Load("../../configs/.env")
	inits.ConnectDatabase("K2-TESTS")
	//Clear collections before tests
	mgm.Coll(new(KBRule)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBObject)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})

	_ebnf = EBNFFactory("../../configs/k2.ebnf")
}

func TestKBRuleValidateIndex(t *testing.T) {
	new(KBRule).validateIndex()
}

func TestRuleFactory(t *testing.T) {
	class, _ := KBClassFactory("Class63f35136f8a767c202687dc8", "", "")
	class.AlterClassAddAttribute("Attribute63f35133f8a767c202687daa", "number", "", []string{}, []string{"User"}, 5, 0)
	ObjectFactoryByClass("Object63f35136f8a767c202687dc8", class)
	sampleRulesOK := [][]string{
		{"Rule01", "unconditionally then create a class named by 'Class63f35136f8a767c202687dc8'"},
		{"Rule02", "unconditionally then alter Class63f35136f8a767c202687dc8 add 'Attribute63f35133f8a767c202687daa' as String from ( User )"},
		{"Rule03", "unconditionally then create an instance of the Class63f35136f8a767c202687dc8 named by 'Object63f35136f8a767c202687dc8'"},
		{"Rule04", "unconditionally then set the Attribute63f35133f8a767c202687daa of the Object63f35136f8a767c202687dc8 = 1243291666028378437"},
	}

	for _, test := range sampleRulesOK {
		time.Sleep(time.Microsecond)
		priority := byte(rand.Intn(100))
		interval := rand.Intn(5000)
		result, err := RuleFactory(test[0], test[1], priority, interval)
		if err != nil {
			t.Errorf("RuleFactory(%v, %v,%v,%v) => %v, %v", test[0], test[1], priority, interval, result, err)
		}
	}

	sampleRulesBad := [][]string{
		{"E01", "unconditionally then set the Attribute63f35133f8a767c202687daa of Object63f35136f8a767c202687dc8 = 1243291666028378437"},
	}
	for _, test := range sampleRulesBad {
		time.Sleep(time.Microsecond)
		priority := byte(rand.Intn(100))
		interval := rand.Intn(5000)
		result, err := RuleFactory(test[0], test[1], priority, interval)
		if err == nil {
			t.Errorf("RuleFactory(%v, %v,%v,%v) => %v, %v", test[0], test[1], priority, interval, result, err)
		}
	}
}

func TestKBRuleClear(t *testing.T) {
	mgm.Coll(new(KBRule)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBObject)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
	t.Log("all clean.")
}
