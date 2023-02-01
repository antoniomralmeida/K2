package models

import "github.com/antoniomralmeida/k2/internal/inits"

type BIN struct {
	tokentype        Tokentype
	literalbin       LiteralBin
	token            string
	class            *KBClass
	attribute        *KBAttribute
	workspace        *KBWorkspace
	objects          []*KBObject
	attributeObjects []*KBAttributeObject
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
