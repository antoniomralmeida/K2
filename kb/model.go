package kb

import (
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KBAttributeType string

const (
	NotDefined KBAttributeType = ""
	KBString   KBAttributeType = "String"
	KBDate     KBAttributeType = "Date"
	KBNumber   KBAttributeType = "Number"
	KBList     KBAttributeType = "List"
)

type KBSource int8

const (
	Empty KBSource = iota
	User
	IOT
	Simulation
	Inference
)

var KBSourceStr = map[string]KBSource{
	"":           Empty,
	"User":       User,
	"IOT":        IOT,
	"Inference":  Inference,
	"Simulation": Simulation,
}

type KBSimulation int8

const (
	Default KBSimulation = iota
	MonteCarlo
	NormalDistribution
	LinearRegression
)

var KBSimulationStr = map[string]KBSimulation{
	"":                   Default,
	"MonteCarlo":         MonteCarlo,
	"NormalDistribution": NormalDistribution,
	"LinearRegression":   LinearRegression,
}

type KnowledgeBased struct {
	mgm.DefaultModel    `json:",inline" bson:",inline"`
	Name                string                          `bson:"name"`
	Classes             []KBClass                       `bson:"-"`
	IdxClasses          map[primitive.ObjectID]*KBClass `bson:"-"`
	Rules               []KBRule                        `bson:"-"`
	Workspaces          []KBWorkspace                   `bson:"-"`
	Objects             []KBObject                      `bson:"-"`
	IdxObjects          map[string]*KBObject            `bson:"-"`
	IdxAttributeObjects map[string]*KBAttributeObject   `bson:"-"`
	ebnf                *ebnf.EBNF                      `bson:"-"`
	stack               []*KBRule                       `bson:"-"`
	mutex               sync.Mutex                      `bson:"-"`
	halt                bool                            `bson:"-"`
}

type KBAttribute struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string          `bson:"name"`
	AType            KBAttributeType `bson:"atype"`
	KeepHistory      int             `bson:"keephistory"`      //Numero de historico a manter, 0- sempre
	ValidityInterval int64           `bson:"validityinterval"` //validade do ultimo valor em microssegudos, 0- sempre
	SimulationID     KBSimulation    `bson:"simulation,omitempty" json:"-"`
	Simulation       string          `bson:"-" json:"simulation"`
	SourcesID        []KBSource      `bson:"sources"`
	Options          []string        `bson:"options,omitempty"`
	Sources          []string        `bson:"-" json:"sources"`
	antecedentRules  []*KBRule       `bson:"-"`
	consequentRules  []*KBRule       `bson:"-"`
}

type KBClass struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string             `bson:"name"`
	Icon             string             `bson:"icon"`
	ParentID         primitive.ObjectID `bson:"parent_id,omitempty"`
	Parent           string             `bson:"-" json:"parent"`
	Attributes       []KBAttribute      `bson:"attributes"`
	ParentClass      *KBClass           `bson:"-"`
}

type BIN struct {
	tokentype        ebnf.Tokentype
	literalbin       models.LiteralBin
	token            string
	class            *KBClass
	attribute        *KBAttribute
	workspace        *KBWorkspace
	objects          []*KBObject
	attributeObjects []*KBAttributeObject
}

type KBRule struct {
	mgm.DefaultModel  `json:",inline" bson:",inline"`
	Rule              string     `bson:"rule"`
	Priority          byte       `bson:"priority"` //0..100
	ExecutionInterval int        `bson:"interval"`
	lastexecution     time.Time  `bson:"-"`
	consequent        int        `bson:"-"`
	inRun             bool       `bson:"-"`
	bkclasses         []*KBClass `bson:"-"`
	bin               []*BIN     `bson:"-"`
}

type KBHistory struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Attribute        primitive.ObjectID `bson:"attribute_id"`
	When             int64              `bson:"when"`
	Value            any                `bson:"value"`
	Trust            float64            `bson:"trust,omitempty"`
	Source           KBSource           `bson:"source"`
}

type KBAttributeObject struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Attribute        primitive.ObjectID `bson:"attribute_id"  json:"AttributeId"`
	KbObject         *KBObject          `bson:"-" json:"-"`
	KbHistory        *KBHistory         `bson:"-" json:"History"`
	KbAttribute      *KBAttribute       `bson:"-"  json:"Attrinute"`
}

type KBObject struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string              `bson:"name"`
	Class            primitive.ObjectID  `bson:"class_id"`
	Attributes       []KBAttributeObject `bson:"attributes"`
	Bkclass          *KBClass            `bson:"-" json:"Class"`
	parsed           bool                `bson:"-"`
}

type KBObjectWS struct {
	Object   primitive.ObjectID `bson:"object_id"`
	Top      int                `bson:"top"`
	Left     int                `bson:"left"`
	KBObject *KBObject          `bson:"-"`
}

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

type DataInput struct {
	Name    string          `json:"name"`
	Atype   KBAttributeType `json:"atype"`
	Options []string        `json:"options"`
}
