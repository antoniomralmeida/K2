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
)

type Token struct {
	id        int
	tokentype Tokentype
	rule_id   int
	rule_jump int
	token     string
	next      []*Token
}

type Statement struct {
	id     int
	name   string
	Tokens []*Token
}

type EBNF struct {
	rules []*Statement
	base  *Token
}

type SYMBOL struct {
	begin string
	end   string
}

type PAIR struct {
	begin int
	end   int
}

var symbols = []SYMBOL{SYMBOL{"=", "."}, SYMBOL{"{", "}"}, SYMBOL{"[", "]"}, SYMBOL{"(", ")"}, SYMBOL{"\"", "\""}, SYMBOL{"'", "'"}}

var TokentypeStr = []string{"", "Reference", "Literal", "Text", "Control", "Jump", "Object", "DynamicReference", "Attribute", "Constant", "Class", "ListType"}
