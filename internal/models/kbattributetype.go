package models

import "strings"

type KBAttributeType int

const (
	NotDefined KBAttributeType = iota
	KBString
	KBDate
	KBNumber
	KBList
)

func KBattributeTypeStr(str string) KBAttributeType {
	return attributeTypeMap[strings.ToLower(str)]
}

var attributeTypeMap = map[string]KBAttributeType{
	"string": KBString,
	"date":   KBDate,
	"number": KBNumber,
	"list":   KBList,
}
