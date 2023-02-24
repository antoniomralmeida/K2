package models

import (
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/kamva/mgm/v3"
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

	EBNFFactory("../../configs/k2.ebnf")
}
