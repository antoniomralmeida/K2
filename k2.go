package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/tests"
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
	tests.Test6(&kbase)
	//tests.Test1(&kbase)

	//time.Sleep(60 * time.Second)
	StartSystem()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(3)
	go kbase.Run(&wg)
	go apikernel.Run(&wg, &kbase)
	go web(&wg)
	wg.Wait()
}

func web(wg *sync.WaitGroup) {
	//WEB SERVER
	switch runtime.GOOS {
	case "windows":
		wd, _ := os.Getwd()
		web := wd + "\\k2web\\k2web.exe"
		b, err := exec.Command("cmd.exe", "/c", "start", web).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(b)
	default:
		log.Fatal("OS not found!" + runtime.GOOS)
	}
	wg.Done()
}
