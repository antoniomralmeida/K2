package models

type EBNF struct {
	Rules []*Statement `json:"rules"`
	Base  *Token       `json:"-"`
}
