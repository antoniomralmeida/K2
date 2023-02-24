package models

import (
	"encoding/json"
	"testing"

	"github.com/kamva/mgm/v3"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestKBObjectValidateIndex(t *testing.T) {
	new(KBObject).ValidateIndex()
}

func TestObjectFactory(t *testing.T) {
	c1 := "Teste" + primitive.NewObjectID().Hex()
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
						t.Errorf("ObjectFactoryByClass(%v,%v) => %v, %v", "Object1", child, obj, err)
					} else {
						if len(obj.Attributes) != 2 {
							t.Errorf("ObjectFactoryByClass(%v,%v) => %v, %v", "Object1", child, obj, err)
						}
					}

					obj2, err := ObjectFactoryByClass("Object1", child)
					if err == nil {
						t.Errorf("ObjectFactoryByClass(%v,%v) => %v, %v", "Object1", child, obj2, err)
					}

					obj3, err := ObjectFactory("Object2", c1)
					if err != nil {
						t.Errorf("ObjectFactory(%v,%v) => %v, %v", "Object2", c1, obj3, err)
					} else {
						if len(obj3.Attributes) != 1 {
							t.Errorf("ObjectFactory(%v,%v) => %v, %v", "Object2", c1, obj3, len(obj3.Attributes))
						}
						j := obj3.String()
						jx := new(KBObject)
						json.Unmarshal([]byte(j), jx)
						if obj3.Name != jx.Name {
							t.Errorf("String() => %v", j)
						}
						obj3.Delete()
					}
					c1 := "Teste" + primitive.NewObjectID().Hex()
					obj4, err := ObjectFactory("Object4", c1)
					if err == nil {
						t.Errorf("ObjectFactory(%v,%v) => %v, %v", "Object4", c1, obj4, err)
					}
				}
			}
		}
	}
}

func TestFindObject(t *testing.T) {
	class, _ := KBClassFactory("Teste"+primitive.NewObjectID().Hex(), "", "")
	class.AlterClassAddAttribute("nome", "string", "", []string{}, []string{"User"}, 5, 0)
	name := "Teste" + primitive.NewObjectID().Hex()
	obj, err := ObjectFactoryByClass(name, class)
	if err == nil {
		obj2 := FindObjectByName(name)
		if obj2 == nil {
			t.Errorf("FindObjectByName(%v) => %v", name, obj2)
		} else {
			if obj.ID != obj2.ID {
				t.Errorf("FindObjectByName(%v) => %v", name, obj2.ID)
			}
		}
		objs, err := FindAllObjects(bson.M{}, "name")
		if err != nil {
			t.Errorf("FindAllObjects(%v, %v) => %v, %v", bson.M{}, "name", objs, err)

		} else {
			if len(objs) < 1 {
				t.Errorf("FindAllObjects(%v, %v) => %v, %v", bson.M{}, "name", len(objs), err)
			}
		}
	}
}

func TestKBObjectClear(t *testing.T) {
	mgm.Coll(new(KBObject)).DeleteMany(mgm.Ctx(), bson.D{{}})
	mgm.Coll(new(KBClass)).DeleteMany(mgm.Ctx(), bson.D{{}})
	t.Log("all clean.")
}
