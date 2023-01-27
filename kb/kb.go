package kb

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/kamva/mgm/v3"
	"github.com/madflojo/tasks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GKB *KnowledgeBased
var scheduler *tasks.Scheduler

func Init() {
	initializers.Log("Init KB", initializers.Info)
	k := KnowledgeBased{}
	GKB = &k
	GKB.FindOne()
	if GKB.Name == "" {
		GKB.Name = "K2 KnowledgeBase System "
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
		GKB.IdxClasses[GKB.Classes[j].ID] = &GKB.Classes[j]
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
				return attrs[i].ID.Hex() < attrs[j].ID.Hex()
			})
			for k, x := range o.Attributes {
				GKB.Objects[j].Attributes[k].KbObject = &GKB.Objects[j]
				//kb.Objects[j].Attributes[k].Kb = kb
				for l, y := range attrs {
					if y.ID == x.Attribute {
						GKB.Objects[j].Attributes[k].KbAttribute = attrs[l]
						break
					}
					if y.ID.Hex() > x.Attribute.Hex() {
						break
					}
				}
				if GKB.Objects[j].Attributes[k].KbAttribute == nil {
					initializers.Log("Attribute not found "+x.Attribute.Hex(), initializers.Fatal)
				}
				GKB.IdxAttributeObjects[o.Name+"."+GKB.Objects[j].Attributes[k].KbAttribute.Name] = &GKB.Objects[j].Attributes[k]

				//Obter ultimo valor
				h := KBHistory{}
				err := h.FindLast(bson.D{{Key: "attribute_id", Value: x.ID}})
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
	collection := mgm.Coll(new(KBWorkspace))
	idx := collection.Indexes()
	ret, err := idx.List(mgm.Ctx())
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		initializers.CreateUniqueIndex(collection, "workspace")
	}
	cursor, err := collection.Find(mgm.Ctx(), bson.D{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(mgm.Ctx(), &GKB.Workspaces)
	return err
}

func FindAllClasses(sort string, cs *[]KBClass) error {
	collection := mgm.Coll(new(KBClass))
	idx := collection.Indexes()
	ret, err := idx.List(context.TODO())
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		initializers.CreateUniqueIndex(collection, "name")
	}
	cursor, err := collection.Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(mgm.Ctx(), cs)
	return err
}

func FindAllObjects(filter bson.M, sort string, os *[]KBObject) error {
	collection := mgm.Coll(new(KBObject))
	idx := collection.Indexes()
	ret, err := idx.List(mgm.Ctx())
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		initializers.CreateUniqueIndex(collection, "name")
	}
	cursor, err := collection.Find(mgm.Ctx(), filter, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(mgm.Ctx(), os)
	return err
}

func FindAllRules(sort string) error {
	collection := mgm.Coll(new(KBRule))
	cursor, err := collection.Find(mgm.Ctx(), bson.M{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	initializers.Log(err, initializers.Fatal)
	err = cursor.All(mgm.Ctx(), &GKB.Rules)
	return err
}

func Run(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start the Scheduler
	scheduler = tasks.New()
	defer scheduler.Stop()

	// Add tasks
	_, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(2 * time.Second),

		TaskFunc: func() error {
			go GKB.RunStackRules()
			return nil
		},
	})
	initializers.Log(err, initializers.Fatal)
	_, err = scheduler.Add(&tasks.Task{
		Interval: time.Duration(60 * time.Second),
		TaskFunc: func() error {
			go GKB.RefreshRules()
			return nil
		},
	})
	initializers.Log(err, initializers.Fatal)

	initializers.Log("K2 KB System started!", initializers.Info)
	if runtime.GOOS == "windows" {
		fmt.Println("K2 KB System started! Press ESC to shutdown")
	}
	for {
		if lib.KeyPress() == 27 || GKB.halt {
			fmt.Printf("Shutdown...")
			scheduler.Stop()
			wg.Done()
			os.Exit(0)
		}
	}

}

func Stop() {
	scheduler.Stop()
}
