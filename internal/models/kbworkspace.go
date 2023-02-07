package models

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Posts            lib.Queue    `bson:"-"`
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
		_kb.Workspaces = append(_kb.Workspaces, w)
		return &w
	} else {
		inits.Log(err, inits.Fatal)
		return nil
	}
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
