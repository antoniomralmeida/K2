package kb

func (me TokenBin) String() string {
	return TokenBinStr[me]
}

func (me TokenBin) Size() int {
	return len(TokenBinStr)
}
