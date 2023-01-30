package kb

import (
	"github.com/antoniomralmeida/k2/cmd/k2/ebnf"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/models"
)

func (b *BIN) GetToken() string {
	return b.token
}

func (b *BIN) GetTokentype() ebnf.Tokentype {
	return b.tokentype
}

func (b *BIN) setTokenBin() {
	if b.GetTokentype() == ebnf.Literal {
		var ok bool
		if b.literalbin, ok = models.LiteralBinStr[b.token]; !ok {
			inits.Log("Literal unknown!"+b.GetToken(), inits.Fatal)
		}
	}
}

func (b *BIN) String() string {
	return "token: " + b.token + ", type:" + b.tokentype.String() + ", bin:" + b.literalbin.String()
}
