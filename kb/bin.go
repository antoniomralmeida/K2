package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/ebnf"
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
			log.Fatal("Literal unknown!", b.GetToken())
		}
	}
}

func (b *BIN) String() string {
	return "token: " + b.token + ", type:" + b.tokentype.String() + ", bin:" + b.literalbin.String()
}
