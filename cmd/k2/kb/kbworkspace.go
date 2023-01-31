package models

import (
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/kamva/mgm/v3"
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
