package tests

import (
	"math/rand"

	"github.com/antoniomralmeida/k2/kb"
)

func Test1(kbase *kb.KnowledgeBased) {
	a := kbase.FindAttributeObjectByName("M01.Potência")
	a.LinearRegression()
}

func Test2(kbase *kb.KnowledgeBased) {
	a := kbase.FindAttributeObjectByName("M01.Potência")
	for i := 0; i < 100; i++ {
		a.SetValue(rand.Float64(), kb.KBSource(kb.Simulation), 50.0)
	}
}
