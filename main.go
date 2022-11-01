package main

import (
	"sync"

	"github.com/antoniomralmeida/k2/db"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
)

var wg sync.WaitGroup = sync.WaitGroup{}

func main() {
	lib.LogInit()
	db.ConnectDB("mongodb://localhost:27017", "K2")
	var kbbase = kb.KnowledgeBase{}
	kbbase.Init("./ebnf/k2.ebnf")
	a := kbbase.FindAttributeObjectByName("M01.Potência")

	a.NormalDistribution()
	wg.Add(10)
	//go kbbase.Run(&wg)
	//go web.Run(&wg, &kbbase)
	wg.Wait()
}
