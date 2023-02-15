package models

import "github.com/antoniomralmeida/k2/internal/inits"

type BIN struct {
	tokentype        Tokentype
	pcNextCommand    int //TODO: para comando sem parametros dinâmicos deve-se saltar para o proxim comando na execução
	literalbin       LiteralBin
	token            string
	class            *KBClass
	newAttributes    []KBAttribute
	attribute        *KBAttribute
	workspace        *KBWorkspace
	objects          []*KBObject          //TODO: Poderia ser dinâmico? Tempo de execução?
	attributeObjects []*KBAttributeObject //TODO: Poderia ser dinâmico? Tempo de execução?
}

func (b *BIN) GetToken() string {
	return b.token
}

func (b *BIN) GetTokentype() Tokentype {
	return b.tokentype
}

func (b *BIN) setTokenBin() {
	if b.GetTokentype() == Literal {
		var ok bool
		if b.literalbin, ok = LiteralBinStr[b.token]; !ok {
			inits.Log("Literal unknown!"+b.GetToken(), inits.Fatal)
		}
	}
}

func (b *BIN) String() string {
	return "token: " + b.token + ", type:" + b.tokentype.String() + ", bin:" + b.literalbin.String()
}
