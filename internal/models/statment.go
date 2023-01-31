package models

type Statement struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Tokens []*Token `json:"tokens"`
}
