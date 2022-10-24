package knowledgebase

import (
	"time"

	"github.com/antoniomralmeida/k2/ebnf"
	"gopkg.in/mgo.v2/bson"
)

type KBAttributeType string

const (
	KBString KBAttributeType = "String"
	KBDate                   = "Date"
	KBNumber                 = "Number"
	KBList                   = "List"
)

type KBSource string

const (
	User       KBSource = "User"
	PLC                 = "PLC"
	History             = "History"
	Simulation          = "Simulation"
)

type KBSimulation string

const (
	Default       KBSimulation = ""
	MonteCarlo                 = "Monte Carlo"
	MovingAverage              = "Moving Average"
	Interpolation              = "interpolation"
)

type TokenBin byte

const (
	b_null TokenBin = iota
	b_open_par
	b_close_par
	b_iqual
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

var TokenBinStr = []string{
	"",
	"(",
	")",
	"=",
	"activate",
	"and",
	"any",
	"change",
	"conclude",
	"create",
	"deactivate",
	"delete",
	"different",
	"equal",
	"focus",
	"for",
	"greater",
	"halt",
	"hide",
	"if",
	"inform",
	"initially",
	"insert",
	"invoke",
	"is",
	"less",
	"move",
	"of",
	"operator",
	"or",
	"remove",
	"rotate",
	"set",
	"show",
	"start",
	"than",
	"that",
	"the",
	"then",
	"to",
	"transfer",
	"unconditionally",
	"when",
	"whenever"}

type KnowledgeBase struct {
	Classes    []KBClass
	Rules      []KBRule
	Workspaces []KBWorkspace
	Objects    []KBObject
	ebnf       *ebnf.EBNF
	stack      []*KBRule
}

type KBAttribute struct {
	Id               bson.ObjectId   `bson:"id,omitempty"`
	Name             string          `bson:"name"`
	AType            KBAttributeType `bson:"atype"`
	Options          []string        `bson:"options,omitempty"`
	Sources          []KBSource      `bson:"sources"`
	KeepHistory      int             `bson:"keephistory"`
	ValidityInterval int             `bson:"validityinterval"`
	Deadline         int             `bson:"deadline"`
	Simulation       KBSimulation    `bson:"simulation,omitempty"`
}

type KBClass struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Name            string        `bson:"name"`
	Icon            string        `bson:"icon"`
	Parent          bson.ObjectId `bson:"parent_id,omitempty"`
	Attributes      []KBAttribute `bson:"attributes"`
	ParentClass     *KBClass      `bson:"-"`
	antecedentRules []*KBRule     `bson:"-"`
	consequentRules []*KBRule     `bson:"-"`
}

type BIN struct {
	Tokentype       ebnf.Tokentype
	typebin         TokenBin
	token           string
	class           *KBClass
	object          *KBObject
	attribute       *KBAttribute
	attributeObject *KBAttributeObject
}

type KBRule struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Rule     string        `bson:"rule"`
	Priority byte          `bson:"priority"` //0..100
	bin      []*BIN        `bson:"-"`
}

type KBHistory struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Attribute bson.ObjectId `bson:"attribute_id"`
	When      time.Time     `bson:"when"`
	Value     string        `bson:"value"`
	Certainty float64       `bson:"certainty,omitempty"`
	Source    KBSource      `bson:"source"`
}

type KBAttributeObject struct {
	Id          bson.ObjectId `bson:"id"`
	Attribute   bson.ObjectId `bson:"attribute_id"`
	KbAttribute *KBAttribute  `bson:"-"`
	KbHistory   *KBHistory    `bson:"-"`
}

type KBObject struct {
	Id         bson.ObjectId       `bson:"_id"`
	Name       string              `bson:"name"`
	Class      bson.ObjectId       `bson:"class_id"`
	Top        int                 `bson:"top"`
	Left       int                 `bson:"left"`
	Attributes []KBAttributeObject `bson:"attributes"`
	Bkclass    *KBClass            `bson:"-"`
}

type KBWorkspace struct {
	Id              bson.ObjectId   `bson:"_id,omitempty"`
	Workspace       string          `bson:"workspace"`
	Top             int             `bson:"top"`
	Left            int             `bson:"left"`
	Width           int             `bson:"width"`
	Height          int             `bson:"height"`
	BackgroundImage string          `bson:"backgroundimage,omitempty"`
	Objects         []bson.ObjectId `bson:"objects"`
	KBObjects       []*KBObject     `bson:"-"`
}
