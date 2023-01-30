package ebnf

func (me Tokentype) String() string {
	return TokentypeStr[me]
}

func (me Tokentype) Size() int {
	return len(TokentypeStr)
}
