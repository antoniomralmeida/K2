package models

func (me LiteralBin) String() string {
	return string(me)
}

func (me LiteralBin) Size() int {
	return len(LiteralBinStr)
}
