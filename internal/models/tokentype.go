package models

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
	Rule
)

var TokentypeStr = []string{"", "Reference", "Literal", "Text", "Control", "Jump", "Object", "DynamicReference", "Attribute", "Constant", "Class", "ListType", "Workspace", "Rule"}

func (me Tokentype) String() string {
	return TokentypeStr[me]
}

func (me Tokentype) Size() int {
	return len(TokentypeStr)
}
