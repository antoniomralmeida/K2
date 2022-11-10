package main

import (
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
)

//var kbase = kb.KnowledgeBased{}

func init() {
	initializers.InitEnvVars()
	initializers.LogInit()
	initializers.ConnectDB()
	kb.Init()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(3)
	go kb.Run(&wg)
	go apikernel.Run(&wg)
	//go web.Run(&wg)
	wg.Wait()
}

func main() {
	//TEST
	//tests.Test1(&kbase)
	//tests.Test2(&kbase)
	//tests.Test6(&kbase)
	//tests.Test1(&kbase)

	//CORE
	StartSystem()
}
