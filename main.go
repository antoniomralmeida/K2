package main

import (
	"time"

	"github.com/antoniomralmeida/k2/db"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/parallel"
	"github.com/antoniomralmeida/k2/web"
)

func main() {
	lib.LogInit()
	db.ConnectDB("mongodb://localhost:27017", "K2")

	kb1 := kb.KnowledgeBase{}
	kb1.Init()

	kb1.ReadEBNF("./ebnf/k2.ebnf")
	kb1.ReadBK()

	t1 := func() error {
		return kb1.Scan()
	}

	t2 := func() error {
		web.Run()
		return nil
	}

	timeout := time.After(10 * time.Second)

	for {
		select {
		case err := <-parallel.Run(t1, t2):
			lib.LogFatal(err)
		case <-timeout:
			break
		}
		time.Sleep(1 * time.Second)
	}

}
