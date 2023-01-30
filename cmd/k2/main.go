package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gookit/color"

	"github.com/antoniomralmeida/k2/cmd/k2/apikernel"
	"github.com/antoniomralmeida/k2/cmd/k2/kb"
	"github.com/antoniomralmeida/k2/cmd/k2olivia/util"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"
)

func init() {
	// Print the K2 ascii text 3D
	k2ASCII := string(util.ReadFile("config/k2.txt"))
	fmt.Println(color.FgLightGreen.Render(k2ASCII))

	msg := fmt.Sprintf("Initializing K2 KB System, version: %v build: %v PID: %v", lib.GetVersion(), lib.GetBuild(), os.Getppid())
	fmt.Println(msg)
	inits.InitEnvVars()
	inits.LogInit("k2log")
	inits.Log(msg, inits.Info)
	inits.InitTelemetry()
	ctx, spanbase := inits.Begin("main-init", nil)
	_, span := inits.Begin("ConnectDB", ctx)
	inits.ConnectDB()
	span.End()
	_, span = inits.Begin("kb.Init", ctx)
	kb.Init()
	span.End()
	_, span = inits.Begin("InitSecurity", ctx)
	models.InitSecurity()
	span.End()
	_, span = inits.Begin("InitOlivia", ctx)
	//inits.InitOlivia()
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
