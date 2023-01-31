package models

type DataInput struct {
	Name    string          `json:"name"`
	Atype   KBAttributeType `json:"atype"`
	Options []string        `json:"options"`
}
