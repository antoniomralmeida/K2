package models

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/madflojo/tasks"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	kb_current  *KnowledgeBased
	scheduler   *tasks.Scheduler
	_classes    []KBClass
	_workspaces []KBWorkspace
	_rules      []KBRule
	_objects    []KBObject
	_ebnf       *EBNF
)

type KnowledgeBased struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string `bson:"name"`
	RestartFlag      bool   `bson:"-"`
}

func KnowledgeBasedFacotory() *KnowledgeBased {
	kb := new(KnowledgeBased)
	kb.findOne()
	if kb.Name == "" {
		kb.Name = "K2 KnowledgeBase System "
	}
	kb.Persist()
	return kb
}

func pauseKB() {
	scheduler.Lock()
}

func resumeKB() {
	if scheduler != nil {
		scheduler.Unlock()
	}
}

func stopKB() {
	if scheduler != nil {
		scheduler.Stop()
	}
}

func restartKB() {
	if scheduler != nil {
		scheduler.Stop()
		InitKB()
		kb_current.RestartFlag = false
	}
}

func (kb *KnowledgeBased) linkObjects(ws *KBWorkspace, obj *KBObject, left int, top int) {
	ows := KBObjectWS{Object: obj.ID, Left: left, Top: top, KBObject: obj}
	ws.Objects = append(ws.Objects, ows)
	inits.Log(ws.Persist(), inits.Fatal)
}

func (kb *KnowledgeBased) updateKB(name string) error {
	kb.Name = name
	return kb.Persist()
}

func (obj *KnowledgeBased) Persist() error {
	return inits.Persist(obj)

}

func (obj *KnowledgeBased) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (kb *KnowledgeBased) findOne() error {
	ret := mgm.Coll(kb).FindOne(mgm.Ctx(), bson.D{})
	ret.Decode(kb)
	return nil
}

func InitKB() {
	inits.Log("Init KB", inits.Info)
	kb_current = KnowledgeBasedFacotory()

	//Check unique index from database collections
	new(KBClass).validateIndex()
	new(KBWorkspace).validateIndex()
	new(KBObject).validateIndex()

	_ebnf = EBNFFactory("./configs/k2.ebnf")

	FindAllClasses("_id")

	_idxClasses := make(map[primitive.ObjectID]*KBClass)
	for _, c := range _classes {
		_idxClasses[c.ID] = &c
	}

	for j, c := range _classes {
		inits.Log("Prepare Class "+c.Name, inits.Info)
		if !c.ParentID.IsZero() {
			pc := _idxClasses[c.ParentID]
			if pc != nil {
				_classes[j].ParentClass = pc
			} else {
				inits.Log("Parent of Class "+c.Name+" not found!", inits.Fatal)
			}
		}
	}

	_objects, _ := FindAllObjects(bson.M{}, "name")
	for j, o := range _objects {
		c := _idxClasses[o.Class]
		if c != nil {
			_objects[j].Bkclass = c
			attrs := c.FindAttributes()
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].ID.Hex() < attrs[j].ID.Hex()
			})
			for k, x := range o.Attributes {
				_objects[j].Attributes[k].KbObject = &_objects[j]
				//kb.Objects[j].Attributes[k].Kb = kb
				for l, y := range attrs {
					if y.ID == x.Attribute {
						_objects[j].Attributes[k].KbAttribute = attrs[l]
						break
					}
					if y.ID.Hex() > x.Attribute.Hex() {
						break
					}
				}
				if _objects[j].Attributes[k].KbAttribute == nil {
					inits.Log("Attribute not found "+x.Attribute.Hex(), inits.Fatal)
				}

				//Last value
				h := KBHistory{}
				err := h.FindLast(bson.D{{Key: "attribute_id", Value: x.ID}})
				if err != nil {
					if err.Error() != "not found" {
						inits.Log(err, inits.Fatal)
					}
				} else {
					_objects[j].Attributes[k].KbHistory = &h
				}
				_objects[j].Attributes[k].Validity()
			}
		} else {
			inits.Log("Class of object "+o.Name+" not found!", inits.Fatal)
		}
	}

	FindAllWorkspaces("name")

	FindAllRules("_id")

	for i := range _rules {
		_, bin, err := parsingRule(_rules[i].Rule)
		inits.Log(err, inits.Fatal)
		linkerRule(&_rules[i], bin)
	}
}

func KBRun(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start the Scheduler
	scheduler = tasks.New()
	defer scheduler.Stop()

	// Add tasks
	_, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(2 * time.Second),

		TaskFunc: func() error {
			go runStackRules()
			return nil
		},
	})
	inits.Log(err, inits.Fatal)
	_, err = scheduler.Add(&tasks.Task{
		Interval: time.Duration(60 * time.Second),
		TaskFunc: func() error {
			go RefreshRules()
			return nil
		},
	})
	inits.Log(err, inits.Fatal)

	inits.Log("K2 KB System started!", inits.Info)
	if runtime.GOOS == "windows" {
		fmt.Println("K2 KB System started! Press ESC to shutdown")
	}
	for {
		if lib.KeyPress() == 27 {
			fmt.Printf("Shutdown...")
			stopKB()
			wg.Done()
			os.Exit(0)
		}
		if kb_current.RestartFlag {
			restartKB()
		}
	}

}
