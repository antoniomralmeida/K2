package models

import (
	"encoding/json"
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
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
}

func FuzzClassFactory(f *testing.F) {

	f.Add("Teste")

	f.Fuzz(func(t *testing.T, a string) {
		result, err := KBClassFactory(a, "", "")
		if (len(a) > 50 || len(a) < 5) && err == nil {
			t.Errorf("ClassFactory(%v,%v,%v) => %v", a, "", "", result)
		}
		if err == nil && result == nil {
			t.Errorf("ClassFactory(%v,%v,%v) => %v", a, "", "", result)
		}
		if result != nil {
			if result.Delete() != nil {
				t.Errorf("ClassFactory() => %v", err)
			}
		}
	})
}

func TestClassFactory(t *testing.T) {
	mt1 := "Teste " + lib.GeneratePassword(25, 0, 5, 5)
	parent, err := KBClassFactory(mt1, "", "")
	if err == nil {
		result, err := KBClassFactory(mt1+"(1)", "", mt1)
		if err != nil {
			t.Errorf("ClassFactory(%v,%v,%v) => %v", mt1+"(1)", "", mt1, result)
		} else {
			result.Delete()
		}

		j := parent.String()

		jx := new(KBClass)
		json.Unmarshal([]byte(j), jx)
		if parent.Name != jx.Name {
			t.Errorf("String() => %v", j)
		}
		parent.Delete()
	}
	mt2 := "Teste " + lib.GeneratePassword(25, 0, 5, 5)
	result, err := KBClassFactory(mt1, "", mt2)
	if err == nil {
		t.Errorf("ClassFactory(%v,%v,%v) => %v", mt1, "", mt2, result)
		result.Delete()
	}
}

func TestAlterClassAddAttribute(t *testing.T) {
	class, err := KBClassFactory("Teste "+lib.GeneratePassword(25, 0, 5, 5), "", "")
	if err == nil {
		result, err := class.AlterClassAddAttribute("nome", "string", "", []string{}, []string{"User"}, 5, 0)
		if err != nil {
			t.Errorf("class.AlterClassAddAttribute(%v,%v,%v,%v,%v,%v,%v) => %v,%v", "nome", "string", "", []string{}, []string{"User"}, 5, 0, result, err)
		} else {
			array := class.FindAttributes()
			if len(array) != 1 {
				t.Errorf("class.FindAttributes() => %v", len(array))
			}

			result = class.FindAttribute("nome")
			if result == nil {
				t.Errorf("class.FindAttribute(%v) => %v", "nome", result)
			}
		}
		result, err = class.AlterClassAddAttribute("X", "string", "", []string{}, []string{"User"}, 5, 0)
		if err == nil {
			t.Errorf("class.AlterClassAddAttribute(%v,%v,%v,%v,%v,%v,%v) => %v,%v", "X", "string", "", []string{}, []string{"User"}, 5, 0, result, err)
		}

		result, err = class.AlterClassAddAttribute("nome", "bool", "", []string{}, []string{"User"}, 5, 0)
		if err == nil {
			t.Errorf("class.AlterClassAddAttribute(%v,%v,%v,%v,%v,%v,%v) => %v,%v", "nome", "bool", "", []string{}, []string{"User"}, 5, 0, result, err)
		}

	}
}

func TestAlterClassAddAttributeParent(t *testing.T) {
	c1 := "Teste " + lib.GeneratePassword(25, 0, 5, 5)
	parent, _ := KBClassFactory(c1, "", "")
	if parent != nil {
		child, _ := KBClassFactoryParent(c1+"=child", "", parent)
		if child != nil {
			_, err := parent.AlterClassAddAttribute("nome", "string", "", []string{}, []string{"User"}, 5, 0)
			if err == nil {
				_, err := child.AlterClassAddAttribute("endereço", "string", "", []string{}, []string{"User"}, 5, 0)
				if err == nil {
					array := child.FindAttributes()
					if len(array) != 2 {
						t.Errorf("class.FindAttributes() => %v", len(array))
					}
				}
			}
		}
	}
}
func TestClear(t *testing.T) {
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
	t.Log("all clean.")
}
