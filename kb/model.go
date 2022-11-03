package kb

import (
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/ebnf"
	"gopkg.in/mgo.v2/bson"
)

type KBAttributeType string

const (
	KBString KBAttributeType = "String"
	KBDate   KBAttributeType = "Date"
	KBNumber KBAttributeType = "Number"
	KBList   KBAttributeType = "List"
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

type LiteralBin byte

const (
	b_null LiteralBin = iota
	b_open_par
	b_close_par
	b_equal_sym
	b_activate
	b_and
	b_any
	b_change
	b_conclude
	b_create
	b_deactivate
	b_delete
	b_different
	b_equal
	b_focus
	b_for
	b_greater
	b_halt
	b_hide
	b_if
	b_inform
	b_initially
	b_insert
	b_invoke
	b_is
	b_less
	b_move
	b_of
	b_operator
	b_or
	b_remove
	b_rotate
	b_set
	b_show
	b_start
	b_than
	b_that
	b_the
	b_then
	b_to
	b_transfer
	b_unconditionally
	b_when
	b_whenever
)

var LiteralBinStr = map[string]LiteralBin{
	"":                b_null,
	"(":               b_open_par,
	")":               b_close_par,
	"=":               b_equal_sym,
	"activate":        b_activate,
	"and":             b_and,
	"any":             b_any,
	"change":          b_change,
	"conclude":        b_conclude,
	"create":          b_create,
	"deactivate":      b_deactivate,
	"delete":          b_delete,
	"different":       b_different,
	"equal":           b_equal,
	"focus":           b_focus,
	"for":             b_for,
	"greater":         b_greater,
	"halt":            b_halt,
	"hide":            b_hide,
	"if":              b_if,
	"inform":          b_inform,
	"initially":       b_initially,
	"insert":          b_insert,
	"invoke":          b_invoke,
	"is":              b_is,
	"less":            b_less,
	"move":            b_move,
	"of":              b_of,
	"operator":        b_operator,
	"or":              b_or,
	"remove":          b_remove,
	"rotate":          b_rotate,
	"set":             b_set,
	"show":            b_show,
	"start":           b_start,
	"than":            b_than,
	"that":            b_than,
	"the":             b_the,
	"then":            b_then,
	"to":              b_to,
	"transfer":        b_transfer,
	"unconditionally": b_unconditionally,
	"when":            b_when,
	"whenever":        b_whenever}

type KnowledgeBased struct {
	Id                  bson.ObjectId                 `bson:"_id,omitempty"`
	Name                string                        `bson:"name"`
	IOTApi              string                        `bson:"iotapi"`
	Classes             []KBClass                     `bson:"-"`
	IdxClasses          map[bson.ObjectId]*KBClass    `bson:"-"`
	Rules               []KBRule                      `bson:"-"`
	Workspaces          []KBWorkspace                 `bson:"-"`
	Objects             []KBObject                    `bson:"-"`
	IdxObjects          map[string]*KBObject          `bson:"-"`
	IdxAttributeObjects map[string]*KBAttributeObject `bson:"-"`
	ebnf                *ebnf.EBNF                    `bson:"-"`
	stack               []*KBRule                     `bson:"-"`
	mutex               sync.Mutex                    `bson:"-"`
}

type KBAttribute struct {
	Id               bson.ObjectId   `bson:"id,omitempty"`
	Name             string          `bson:"name"`
	AType            KBAttributeType `bson:"atype"`
	Options          []string        `bson:"options,omitempty"`
	SourcesID        []KBSource      `bson:"sources"`
	Sources          []string        `bson:"-" json:"sources"`
	KeepHistory      int64           `bson:"keephistory"`
	ValidityInterval int64           `bson:"validityinterval"`
	Deadline         int64           `bson:"deadline"`
	SimulationID     KBSimulation    `bson:"simulation,omitempty" json:"-"`
	Simulation       string          `bson:"-" json:"simulation"`
	antecedentRules  []*KBRule       `bson:"-"`
	consequentRules  []*KBRule       `bson:"-"`
}

type KBClass struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Icon        string        `bson:"icon"`
	ParentID    bson.ObjectId `bson:"parent_id,omitempty"`
	Parent      string        `bson:"-" json:"parent"`
	Attributes  []KBAttribute `bson:"attributes"`
	ParentClass *KBClass      `bson:"-"`
}

type BIN struct {
	tokentype        ebnf.Tokentype
	literalbin       LiteralBin
	token            string
	class            *KBClass
	attribute        *KBAttribute
	objects          []*KBObject
	attributeObjects []*KBAttributeObject
}

type KBRule struct {
	Id                bson.ObjectId `bson:"_id,omitempty"`
	Rule              string        `bson:"rule"`
	Priority          byte          `bson:"priority"` //0..100
	ExecutionInterval int           `bson:"interval"`
	bin               []*BIN        `bson:"-"`
	lastexecution     time.Time     `bson:"-"`
	bkclasses         []*KBClass    `bson:"-"`
	consequent        int           `bson:"-"`
}

type KBHistory struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Attribute bson.ObjectId `bson:"attribute_id"`
	When      int64         `bson:"when"`
	Value     any           `bson:"value"`
	Trust     float32       `bson:"trust,omitempty"`
	Source    KBSource      `bson:"source"`
}

type KBAttributeObject struct {
	Id          bson.ObjectId   `bson:"id"`
	Attribute   bson.ObjectId   `bson:"attribute_id"  json:"AttributeId"`
	KbObject    *KBObject       `bson:"-" json:"-"`
	KbHistory   *KBHistory      `bson:"-" json:"History"`
	KbAttribute *KBAttribute    `bson:"-"  json:"Attrinute"`
	Kb          *KnowledgeBased `bson:"-"  json:"-"`
}

type KBObject struct {
	Id         bson.ObjectId       `bson:"_id"`
	Name       string              `bson:"name"`
	Class      bson.ObjectId       `bson:"class_id"`
	Attributes []KBAttributeObject `bson:"attributes"`
	Bkclass    *KBClass            `bson:"-" json:"Class"`
	parsed     bool                `bson:"-"`
}

type KBObjectWS struct {
	Object   bson.ObjectId `bson:"object_id"`
	Top      int           `bson:"top"`
	Left     int           `bson:"left"`
	KBObject *KBObject     `bson:"-"`
}

type KBWorkspace struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Workspace       string        `bson:"workspace"`
	Top             int           `bson:"top"`
	Left            int           `bson:"left"`
	Width           int           `bson:"width"`
	Height          int           `bson:"height"`
	BackgroundImage string        `bson:"backgroundimage,omitempty"`
	Objects         []KBObjectWS  `bson:"objects"`
}

type DataInput struct {
	Id      bson.ObjectId   `json:"id"`
	Name    string          `json:"name"`
	Atype   KBAttributeType `json:"atype"`
	Options []string        `json:"options"`
}
