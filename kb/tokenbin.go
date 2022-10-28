package kb

func (me TokenBin) String() string {
	return string(me)
}

func (me TokenBin) Size() int {
	return len(TokenBinStr)
}
