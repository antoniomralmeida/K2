package kb

import (
	"context"
	"sort"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GKB *KnowledgeBased

func Init() {
	initializers.Log("Init KB", initializers.Info)
	k := KnowledgeBased{}
	GKB = &k
	GKB.FindOne()
	if GKB.Name == "" {
		GKB.Name = "K2 System KB"
	}
	GKB.Persist()
	GKB.IdxClasses = make(map[primitive.ObjectID]*KBClass)
	GKB.IdxObjects = make(map[string]*KBObject)
	GKB.IdxAttributeObjects = make(map[string]*KBAttributeObject)

	ebnf := ebnf.EBNF{}
	GKB.ebnf = &ebnf
	GKB.ebnf.ReadToken("./config/k2.ebnf")

	FindAllClasses("_id", &GKB.Classes)
	for j := range GKB.Classes {
		GKB.IdxClasses[GKB.Classes[j].Id] = &GKB.Classes[j]
	}

	for j, c := range GKB.Classes {
		initializers.Log("Prepare Class "+c.Name, initializers.Info)
		if !c.ParentID.IsZero() {
			pc := GKB.IdxClasses[c.ParentID]
			if pc != nil {
				GKB.Classes[j].ParentClass = pc
			} else {
				initializers.Log("Parent of Class "+c.Name+" not found!", initializers.Fatal)
			}
		}
	}

	FindAllObjects(bson.M{}, "name", &GKB.Objects)
	for j, o := range GKB.Objects {
		GKB.IdxObjects[o.Name] = &GKB.Objects[j]
		c := GKB.IdxClasses[o.Class]
		if c != nil {
			GKB.Objects[j].Bkclass = c
			attrs := GKB.FindAttributes(c)
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].Id.Hex() < attrs[j].Id.Hex()
			})
			for k, x := range o.Attributes {
				GKB.Objects[j].Attributes[k].KbObject = &GKB.Objects[j]
				//kb.Objects[j].Attributes[k].Kb = kb
				for l, y := range attrs {
					if y.Id == x.Attribute {
						GKB.Objects[j].Attributes[k].KbAttribute = attrs[l]
						break
					}
					if y.Id.Hex() > x.Attribute.Hex() {
						break
					}
				}
				if GKB.Objects[j].Attributes[k].KbAttribute == nil {
					initializers.Log("Attribute not found "+x.Attribute.Hex(), initializers.Fatal)
				}
				GKB.IdxAttributeObjects[o.Name+"."+GKB.Objects[j].Attributes[k].KbAttribute.Name] = &GKB.Objects[j].Attributes[k]

				//Obter ultimo valor
				h := KBHistory{}
				err := h.FindLast(bson.D{{Key: "attribute_id", Value: x.Id}})
				if err != nil {
					if err.Error() != "not found" {
						initializers.Log(err, initializers.Fatal)
					}
				} else {
					GKB.Objects[j].Attributes[k].KbHistory = &h
				}
				GKB.Objects[j].Attributes[k].Validity()
			}
		} else {
			initializers.Log("Class of object "+o.Name+" not found!", initializers.Fatal)
		}
	}

	FindAllWorkspaces("name")

	FindAllRules("_id")

	for i := range GKB.Rules {
		_, bin, err := GKB.ParsingCommand(GKB.Rules[i].Rule)
		initializers.Log(err, initializers.Fatal)
		GKB.linkerRule(&GKB.Rules[i], bin)
	}
}

func FindAllWorkspaces(sort string) error {
	ctx, collection := initializers.GetCollection("KBWorkspace")
	idx := collection.Indexes()
	ret, err := idx.List(ctx)
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		_, err = idx.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"workspace": 1}, Options: options.Index().SetUnique(true)})
		initializers.Log(err, initializers.Fatal)
	}
	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, &GKB.Workspaces)
	return err
}

func FindAllClasses(sort string, cs *[]KBClass) error {
	ctx, collection := initializers.GetCollection("KBClass")
	idx := collection.Indexes()
	ret, err := idx.List(context.TODO())
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		_, err = idx.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)})
		initializers.Log(err, initializers.Fatal)
	}
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, cs)
	return err
}

func FindAllObjects(filter bson.M, sort string, os *[]KBObject) error {
	ctx, collection := initializers.GetCollection("KBObject")
	idx := collection.Indexes()
	ret, err := idx.List(ctx)
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		_, err = idx.CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)})
		initializers.Log(err, initializers.Fatal)
	}
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, os)
	return err
}

func FindAllRules(sort string) error {
	ctx, collection := initializers.GetCollection("KBRule")
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(ctx, &GKB.Rules)
	return err
}
