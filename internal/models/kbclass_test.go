package models

import (
	"encoding/json"
	"testing"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load("../../configs/.env")
	inits.ConnectDB()
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
