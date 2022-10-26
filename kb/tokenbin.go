package kb

import "fmt"

func (me TokenBin) String() string {
	return fmt.Sprintf("%v", me)
}

func (me TokenBin) Size() int {
	return len(TokenBinStr)
}
