package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/web"
	"github.com/subosito/gotenv"
)

var kbase = kb.KnowledgeBase{}

func init() {
	gotenv.Load()
	initializers.ConnectDB()
	initializers.LogInit()
	kbase.Init()
}

func main() {
	//TEST
	//tests()

	StartSystem()
}

func StartSystem() {
	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	tasks, _ := strconv.Atoi(os.Getenv("GOTASKS"))
	wg.Add(tasks)
	go kbase.Run(&wg)
	go web.Run(&wg, &kbase)
	wg.Wait()
}

func Tests() {
	a := kbase.FindAttributeObjectByName("M01.Potência")
	a.NormalDistribution()
}
