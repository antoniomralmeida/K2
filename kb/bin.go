package kb

import (
	"github.com/antoniomralmeida/k2/ebnf"
	"github.com/antoniomralmeida/k2/initializers"
)

func (b *BIN) GetToken() string {
	return b.token
}

func (b *BIN) GetTokentype() ebnf.Tokentype {
	return b.tokentype
}

func (b *BIN) setTokenBin() {
	if b.GetTokentype() == ebnf.Literal {
		b.literalbin = LiteralBinStr[b.token]
		if b.literalbin == b_null {
			initializers.Log("Literal unknown!"+b.GetToken(), initializers.Fatal)
		}
	}
}

func (b *BIN) String() string {
	return "token: " + b.token + ", type:" + b.tokentype.String() + ", bin:" + b.literalbin.String()
}
