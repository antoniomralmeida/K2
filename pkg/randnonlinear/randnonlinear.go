package randnonlinear

import (
	"math/rand"
	"time"
)

func IntnNL(max int, weights func(int) int) int {
	set := []int{}
	for i := 0; i < max; i++ {
		w := weights(i)
		for w > 0 {
			set = append(set, i)
			w--
		}
	}
	time.Sleep(time.Millisecond)
	r := rand.Intn(len(set))
	return set[r]
}
