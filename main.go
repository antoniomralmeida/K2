package main

import (
	"sync"

	"github.com/antoniomralmeida/k2/db"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/web"
)

var wg sync.WaitGroup = sync.WaitGroup{}

func main() {
	lib.LogInit()
	db.ConnectDB("mongodb://localhost:27017", "K2")

	var kbbase = kb.KnowledgeBase{}
	kbbase.Init("./ebnf/k2.ebnf")
	wg.Add(10)
	//go kbbase.Run(&wg)
	go web.Run(&wg, &kbbase)
	wg.Wait()
}
