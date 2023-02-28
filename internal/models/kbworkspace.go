package models

import (
	"context"
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/pkg/queue"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KBWorkspace struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Workspace        string       `bson:"workspace"`
	Top              int          `bson:"top"`
	Left             int          `bson:"left"`
	Width            int          `bson:"width"`
	Height           int          `bson:"height"`
	BackgroundImage  string       `bson:"backgroundimage,omitempty"`
	Objects          []KBObjectWS `bson:"objects"`
	Posts            queue.Queue  `bson:"-"`
}

func WorkspaceFactory(name string, image string) *KBWorkspace {
	copy, err := lib.LoadImage(image)
	if err != nil {
		inits.Log(err, inits.Error)
		return nil
	}
	w := KBWorkspace{Workspace: name, BackgroundImage: copy}
	err = w.Persist()
	if err == nil {
		_workspaces = append(_workspaces, w)
		return &w
	} else {
		inits.Log(err, inits.Fatal)
		return nil
	}
}

func (obj *KBWorkspace) ValidateIndex() error {
	cur, err := mgm.Coll(obj).Indexes().List(mgm.Ctx())
	inits.Log(err, inits.Error)
	var result []bson.M
	err = cur.All(context.TODO(), &result)
	if len(result) == 1 {
		inits.CreateUniqueIndex(mgm.Coll(obj), "workspace")
	}
	return err
}

func (obj *KBWorkspace) Persist() error {
	return inits.Persist(obj)

}

func (obj *KBWorkspace) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (w *KBWorkspace) String() string {
	j, err := json.MarshalIndent(*w, "", "\t")
	inits.Log(err, inits.Error)
	return string(j)
}

func (w *KBWorkspace) AddObject(obj *KBObject, left, top int) {
	ows := new(KBObjectWS)
	ows.KBObject = obj
	ows.Object = obj.ID
	ows.Left = left
	ows.Top = top
	w.Objects = append(w.Objects, *ows)
	w.Persist()
}

func FindAllWorkspaces(sort string) error {
	cursor, err := mgm.Coll(new(KBWorkspace)).Find(mgm.Ctx(), bson.D{}, options.Find().SetSort(bson.D{{Key: sort, Value: 1}}))
	inits.Log(err, inits.Fatal)
	err = cursor.All(mgm.Ctx(), &_workspaces)
	return err
}

func FindWorkspaceByName(name string) *KBWorkspace {
	for i := range _workspaces {
		if _workspaces[i].Workspace == name {
			return &_workspaces[i]
		}
	}
	inits.Log("Workspace not found!", inits.Error)
	return nil
}

func KBWorkspacesJson() (string, error) {
	wks := []KBWorkspace{}

	if err := mgm.Coll(new(KBWorkspace)).SimpleFind(&wks, bson.D{{}}); err != nil {
		return "", err
	}
	ret := []WorkspaceInfo{}
	for _, w := range wks {
		ret = append(ret, WorkspaceInfo{Workspace: w.Workspace, BackgroundImage: w.BackgroundImage})
	}
	json, err := json.Marshal(ret)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func KBGetWorkspacesFromObject(o *KBObject) (ret []*KBWorkspace) {
	//TODO: From mongoDB
	for i := range _workspaces {
		for j := range _workspaces[i].Objects {
			if _workspaces[i].Objects[j].KBObject == o {
				ret = append(ret, &_workspaces[i])
			}
		}
	}
	return
}
