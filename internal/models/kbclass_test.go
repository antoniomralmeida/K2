package models

import (
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
		result, err := ClassFactory(a, "", "")
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
	parent, err := ClassFactory(mt1, "", "")
	if err == nil {
		result, err := ClassFactory(mt1+"(1)", "", mt1)
		if err != nil {
			t.Errorf("ClassFactory(%v,%v,%v) => %v", mt1+"(1)", "", mt1, result)
		} else {
			result.Delete()
		}
		parent.Delete()
	}
	mt2 := "Teste " + lib.GeneratePassword(25, 0, 5, 5)
	result, err := ClassFactory(mt1, "", mt2)
	if err == nil {
		t.Errorf("ClassFactory(%v,%v,%v) => %v", mt1+"(1)", "", mt1, result)
		result.Delete()
	}
}
