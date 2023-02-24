package models

import (
	"encoding/json"
	"testing"

	"github.com/kamva/mgm/v3"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestKBClassValidateIndex(t *testing.T) {
	new(KBClass).ValidateIndex()
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
	KBClassFactory("Class63f35136f8a767c202687dc8", "", "")
	mt1 := "Teste" + primitive.NewObjectID().Hex()
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
	mt2 := "Teste" + primitive.NewObjectID().Hex()
	result, err := KBClassFactory(mt1, "", mt2)
	if err == nil {
		t.Errorf("ClassFactory(%v,%v,%v) => %v", mt1, "", mt2, result)
		result.Delete()
	}
}

func TestAlterClassAddAttribute(t *testing.T) {
	class, err := KBClassFactory("Teste"+primitive.NewObjectID().Hex(), "", "")
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
	c1 := "Teste" + primitive.NewObjectID().Hex()
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

func TestFindAllClasses(t *testing.T) {
	_, err := KBClassFactory("Teste"+primitive.NewObjectID().Hex(), "", "")
	if err == nil {
		cls, err := FindAllClasses("name")
		if err != nil {
			t.Errorf("FindAllClasses(%v) => %v", "name", err)
		}
		if len(*cls) < 1 {
			t.Errorf("FindAllClasses(%v) => %v", "name", len(*cls))
		}
	}
}

func TestKBClassCopy(t *testing.T) {
	parent, _ := KBClassFactory("Teste"+primitive.NewObjectID().Hex(), "teste.jpg", "")
	c1 := "Teste" + primitive.NewObjectID().Hex()
	cl, err := KBClassFactoryParent(c1, "teste2.jpg", parent)
	if err == nil {
		_, err := cl.AlterClassAddAttribute("nome", "string", "", []string{}, []string{"User"}, 5, 0)
		if err == nil {
			cl2, err := KBClassCopy(c1+"(copy)", cl)
			if err != nil {
				t.Errorf("KBClassCopy(%v,%v) => %v,%v", c1+"(copy)", cl.ID, cl2, err)
			}
			if cl2 == nil {
				t.Errorf("KBClassCopy(%v,%v) => %v,%v", c1+"(copy)", cl.ID, cl2, err)
			} else {
				if cl.ParentID != cl2.ParentID {
					t.Errorf("KBClassCopy(%v,%v) => %v,%v", c1+"(copy)", cl.ID, cl, cl2)
				}
				if cl.Icon != cl2.Icon {
					t.Errorf("KBClassCopy(%v,%v) => %v,%v", c1+"(copy)", cl.ID, cl, cl2)
				}
			}
		}
	}
}

func TestKBClassClear(t *testing.T) {
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
	t.Log("all clean.")
}
