package cartesian

//generates a cartesian product on datasets

type item struct {
	idx int
	max int
}

type Cartesian struct {
	itens   map[string]*item
	keys    []string
	current int
}

func (c *Cartesian) AddItem(key string, max int) {
	if c.itens == nil {
		c.itens = make(map[string]*item)
	}
	c.itens[key] = &item{max: max}
	c.keys = []string{}
	c.current = 0
	for key := range c.itens {
		c.keys = append(c.keys, key)
	}
}

func (c *Cartesian) next(i int) bool {
	if i >= len(c.keys) {
		return false
	}
	if c.itens[c.keys[i]].idx < c.itens[c.keys[i]].max {
		c.itens[c.keys[i]].idx = c.itens[c.keys[i]].idx + 1
		return true
	} else if i < len(c.keys) {
		for k := i; k <= i; k++ {
			c.itens[c.keys[k]].idx = 0
		}
		return c.next(i + 1)
	}
	return false
}

func (c *Cartesian) GetCombination() (end bool, idxs map[string]int) {
	idxs = make(map[string]int)
	for key := range c.itens {
		idxs[key] = c.itens[key].idx
	}
	end = c.next(c.current)
	return
}
