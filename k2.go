package main

import (
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/web"
)

var kbase = kb.KnowledgeBased{}

func init() {
	initializers.InitEnvVars()
	initializers.LogInit()
	initializers.ConnectDB()
	kbase.Init()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(3)
	go kbase.Run(&wg)
	go apikernel.Run(&wg, &kbase)
	go web.Run(&wg)
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
