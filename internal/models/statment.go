package models

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
)

type Statement struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Tokens []*Token `json:"tokens"`
}

func (s *Statement) String() string {
	ret, err := json.MarshalIndent(s, "", "    ")
	inits.Log(err, inits.Error)
	return string(ret)
}
