package ebnf

import "strconv"

func (t *Token) GetToken() string {
	return t.Token
}

func (t *Token) GetTokentype() Tokentype {
	return t.Tokentype
}

func (t *Token) GetNexts() []*Token {
	return t.Nexts
}

func (t *Token) String() string {
	return "#" + strconv.Itoa(t.Id) + ", token: " + t.Token + ", type:" + t.Tokentype.String()
}
