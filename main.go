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

	kb1 := kb.KnowledgeBase{}
	kb1.Init("./ebnf/k2.ebnf")
	wg.Add(10)
	go kb1.Run(&wg)
	go web.Run(&wg)
	wg.Wait()
}
