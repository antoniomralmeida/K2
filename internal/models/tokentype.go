package models

type TokenType byte

const (
	Null TokenType = iota
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

func (me TokenType) String() string {
	return TokentypeStr[me]
}

func (me TokenType) Size() int {
	return len(TokentypeStr)
}
