package models

type KBAttributeType string

const (
	NotDefined KBAttributeType = ""
	KBString   KBAttributeType = "String"
	KBDate     KBAttributeType = "Date"
	KBNumber   KBAttributeType = "Number"
	KBList     KBAttributeType = "List"
)
