package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/version"
	"github.com/antoniomralmeida/k2/web"
)

func init() {
	msg := fmt.Sprintf("Initializing K2 system version: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
	fmt.Println(msg)
	initializers.InitEnvVars()
	initializers.LogInit("k2log")
	initializers.Log(msg, initializers.Info)
	initializers.ConnectDB()
	kb.Init()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(3)
	go kb.Run(&wg)
	go apikernel.Run(&wg)
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
