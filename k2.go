package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/k2olivia/util"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/models"
	"github.com/gookit/color"

	"github.com/antoniomralmeida/k2/version"
)

func init() {
	// Print the K2 ascii text 3D
	k2ASCII := string(util.ReadFile("config/k2.txt"))
	fmt.Println(color.FgLightGreen.Render(k2ASCII))

	msg := fmt.Sprintf("Initializing K2 KB System, version: %v build: %v PID: %v", version.GetVersion(), version.GetBuild(), os.Getppid())
	fmt.Println(msg)
	initializers.InitEnvVars()
	initializers.LogInit("k2log")
	initializers.Log(msg, initializers.Info)
	initializers.InitTelemetry()
	ctx, spanbase := initializers.Begin("main-init", nil)
	_, span := initializers.Begin("ConnectDB", ctx)
	initializers.ConnectDB()
	span.End()
	_, span = initializers.Begin("kb.Init", ctx)
	kb.Init()
	span.End()
	_, span = initializers.Begin("InitSecurity", ctx)
	models.InitSecurity()
	span.End()
	_, span = initializers.Begin("InitOlivia", ctx)
	//initializers.InitOlivia()
	span.End()
	spanbase.End()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(5)
	go apikernel.Run(&wg)
	go kb.Run(&wg)
	go lib.Openbrowser("http://localhost" + os.Getenv("HTTPPORT"))
	wg.Wait()
}

func main() {
	//TEST

	//test.Test7()

	//CORE
	StartSystem()
}
