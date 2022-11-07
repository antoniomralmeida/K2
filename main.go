package main

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/tests"
	"github.com/antoniomralmeida/k2/web"
)

var kbase = kb.KnowledgeBased{}

func init() {
	initializers.InitEnvVars()
	initializers.LogInit()
	initializers.ConnectDB()
	kbase.Init()
}

func main() {
	//TEST
	//tests.Test1(&kbase)
	//tests.Test2(&kbase)
	tests.Test3(&kbase)
	//tests.Test1(&kbase)

	time.Sleep(60 * time.Second)
	//StartSystem()
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
