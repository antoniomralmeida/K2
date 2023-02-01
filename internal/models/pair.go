package models

type PAIR struct {
	begin int
	end   int
}

func findPair(p []PAIR, i int) int {
	var ret = 0
	for k, x := range p {
		if x.begin <= i && x.end >= i && (p[k].begin > p[ret].begin || p[k].end < p[ret].end) {
			ret = k
		}
	}
	return ret
}
