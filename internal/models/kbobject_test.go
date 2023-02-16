package models

import (
	"testing"

	"github.com/kamva/mgm/v3"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"

	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	gotenv.Load("../../configs/.env")
	inits.ConnectDatabase("K2-TESTS")
	//Clear collection kb_class before tests
	mgm.Coll(new(KBObject)).DeleteMany(mgm.Ctx(), bson.D{{}})
}

func TestKBObjectValidateIndex(t *testing.T) {
	new(KBObject).validateIndex()
}

func TestObjectFactory(t *testing.T) {
	c1 := "Teste " + lib.GeneratePassword(25, 0, 5, 5)
	parent, _ := KBClassFactory(c1, "", "")
	if parent != nil {
		child, _ := KBClassFactoryParent(c1+"=child", "", parent)
		if child != nil {
			_, err := parent.AlterClassAddAttribute("nome", "string", "", []string{}, []string{"User"}, 5, 0)
			if err == nil {
				_, err := child.AlterClassAddAttribute("endereço", "string", "", []string{}, []string{"User"}, 5, 0)
				if err == nil {
					obj, err := ObjectFactoryByClass("Object1", child)
					if err != nil {
						t.Errorf("models.ObjectFactoryByClass(%v,%v) => %v, %v", "Object1", child, obj, err)
					} else {
						if len(obj.Attributes) != 2 {
							t.Errorf("models.ObjectFactoryByClass(%v,%v) => %v, %v", "Object1", child, obj, err)
						}
					}

					obj2, err := ObjectFactoryByClass("Object1", child)
					if err == nil {
						t.Errorf("models.ObjectFactoryByClass(%v,%v) => %v, %v", "Object1", child, obj2, err)
					}
				}
			}
		}
	}
}

func TestKBObjectClear(t *testing.T) {
	mgm.Coll(new(KBObject)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
	t.Log("all clean.")
}
