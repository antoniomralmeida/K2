package ebnf

type Tokentype byte

const (
	Null Tokentype = iota
	Reference
	Literal
	Text
	Control
	Jump
	Object
	DynamicReference
	Attribute
	Constant
	Class
	ListType
	Workspace
)

type Token struct {
	Id        int       `json:"id"`
	Tokentype Tokentype `json:"tokentype"`
	Rule_id   int       `json:"rule_id"`
	Rule_jump int       `json:"rule_jump"`
	Token     string    `json:"token"`
	Nexts     []*Token  `json:"-"`
}

type Statement struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Tokens []*Token `json:"tokens"`
}

type EBNF struct {
	Rules []*Statement `json:"rules"`
	Base  *Token       `json:"-"`
}

type SYMBOL struct {
	begin string
	end   string
}

type PAIR struct {
	begin int
	end   int
}

var symbols = []SYMBOL{{"=", "."}, {"{", "}"}, {"[", "]"}, {"(", ")"}, {"\"", "\""}, {"'", "'"}}

var TokentypeStr = []string{"", "Reference", "Literal", "Text", "Control", "Jump", "Object", "DynamicReference", "Attribute", "Constant", "Class", "ListType", "Workspace"}
