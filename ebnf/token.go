package ebnf

import "strconv"

func (t *Token) GetToken() string {
	return t.token
}

func (t *Token) GetTokentype() Tokentype {
	return t.tokentype
}

func (t *Token) GetNexts() []*Token {
	return t.next
}

func (t *Token) String() string {
	return "#" + strconv.Itoa(t.id) + ", token: " + t.token + ", type:" + t.tokentype.String()
}
