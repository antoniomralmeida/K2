package kb

import (
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/lib"
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

type LiteralBin byte

const (
	b_a LiteralBin = iota
	b_activate
	b_an
	b_and
	b_any
	b_by
	b_change
	b_cloning
	b_close_par
	b_conclude
	b_create
	b_deactivate
	b_delete
	b_different
	b_equal
	b_equal_sym
	b_focus
	b_for
	b_greater
	b_halt
	b_hide
	b_if
	b_inform
	b_initially
	b_insert
	b_instance
	b_invoke
	b_is
	b_less
	b_move
	b_named
	b_of
	b_open_par
	b_operator
	b_or
	b_parent
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
	b_whose
)

var LiteralBinStr = map[string]LiteralBin{
	"a":               b_a,
	"activate":        b_activate,
	"an":              b_an,
	"and":             b_and,
	"any":             b_any,
	"by":              b_by,
	"change":          b_change,
	"cloning":         b_cloning,
	")":               b_close_par,
	"conclude":        b_conclude,
	"create":          b_create,
	"deactivate":      b_deactivate,
	"delete":          b_delete,
	"different":       b_different,
	"equal":           b_equal,
	"=":               b_equal_sym,
	"focus":           b_focus,
	"for":             b_for,
	"greater":         b_greater,
	"halt":            b_halt,
	"hide":            b_hide,
	"if":              b_if,
	"inform":          b_inform,
	"initially":       b_initially,
	"insert":          b_insert,
	"instance":        b_instance,
	"invoke":          b_invoke,
	"is":              b_is,
	"less":            b_less,
	"move":            b_move,
	"named":           b_named,
	"of":              b_of,
	"(":               b_open_par,
	"operator":        b_operator,
	"or":              b_or,
	"parent":          b_parent,
	"remove":          b_remove,
	"rotate":          b_rotate,
	"set":             b_set,
	"show":            b_show,
	"start":           b_start,
	"than":            b_than,
	"that":            b_that,
	"the":             b_the,
	"then":            b_then,
	"to":              b_to,
	"transfer":        b_transfer,
	"unconditionally": b_unconditionally,
	"when":            b_when,
	"whenever":        b_whenever,
	"whose":           b_whose,
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
	literalbin       LiteralBin
	token            string
	class            *KBClass
	attribute        *KBAttribute
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

type Pipe struct {
	id     primitive.ObjectID `json:"_id"`
	avg    float64            `json:"avg"`
	stdDev float64            `json:"stdDev"`
	trust  float64            `json:"trust"`
}
