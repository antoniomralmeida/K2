package kb

import (
	"log"

	"github.com/antoniomralmeida/k2/ebnf"
)

func (b *BIN) findTokenBin(i byte, j byte) TokenBin {
	if j >= i {
		avg := (i + j) / 2
		tb := TokenBin(avg)
		if b.tokentype.String() == tb.String() {
			return tb
		} else if b.tokentype.String() >= tb.String() {
			return b.findTokenBin(avg+1, j)
		} else {
			return b.findTokenBin(i, avg-1)
		}
	}
	return TokenBin(0)
}

func (b *BIN) GetToken() string {
	return b.token
}

func (b *BIN) GetTokentype() ebnf.Tokentype {
	return b.tokentype
}

func (b *BIN) setTokenBin() {
	if b.GetTokentype() == ebnf.Literal {
		b.typebin = b.findTokenBin(0, byte(b.typebin.Size()-1))
		if b.typebin == b_null {
			log.Fatal("Literal unknown!", b.GetToken())
		}
	}
}

func (b *BIN) String() string {
	return "token: " + b.token + ", type:" + b.tokentype.String() + ", bin:" + b.typebin.String()
}
