package models

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
)

type BIN struct {
	TokenType TokenType `json:"tokentype"`
	//pcNextCommand int        `json:"-"`
	LiteralBin LiteralBin `json:"literalbin"`
	Token      string     `json:"token"`
	class      *KBClass   `json:"-"`
	//newAttributes    []KBAttribute        `json:"-"`
	attribute        *KBAttribute         `json:"-"`
	workspace        *KBWorkspace         `json:"-"`
	objects          []*KBObject          `json:"-"` //TODO: Poderia ser dinâmico? Tempo de execução?
	attributeObjects []*KBAttributeObject `json:"-"` //TODO: Poderia ser dinâmico? Tempo de execução?
}

func (b *BIN) GetToken() string {
	return b.Token
}

func (b *BIN) GetTokentype() TokenType {
	return b.TokenType
}

func (b *BIN) CheckLiteralBin() error {
	if b.GetTokentype() == Literal {
		var ok bool
		if b.LiteralBin, ok = LiteralBinStr[b.Token]; !ok {
			return lib.LiteralNotFoundError
		}
	}
	return nil
}

func (b *BIN) String() string {
	j, err := json.Marshal(*b)
	inits.Log(err, inits.Error)
	return string(j)
}
