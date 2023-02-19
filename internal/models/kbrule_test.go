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

func FuzzRuleFactory(f *testing.F) {

	f.Add("initially Rule")

	f.Fuzz(func(t *testing.T, a string) {
		sampleRule := _ebnf.GrammarSample()
		time.Sleep(time.Microsecond)
		priority := byte(rand.Intn(100))
		interval := rand.Intn(5000)
		result, err := RuleFactory(sampleRule, priority, interval)
		if err != nil {
			t.Errorf("RuleFactory(%v) => %v, %v", sampleRule, result, err)
		}
	})
}

func TestRuleFactory(t *testing.T) {
	//r1 := "Teste " + lib.GeneratePassword(25, 0, 5, 5)

}

func TestKBRuleClear(t *testing.T) {
	mgm.Coll(new(KBRule)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBObject)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
	t.Log("all clean.")
}
