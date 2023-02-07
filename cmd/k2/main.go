package main

import (
	"fmt"

	"os"
	"sync"

	"github.com/gookit/color"

	"github.com/antoniomralmeida/k2/internal/apikernel"
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"
	"github.com/antoniomralmeida/k2/pkg/version"
)

func init() {
	// Print the K2 ascii text 3D
	k2ASCII, _ := os.ReadFile("./configs/k2.txt")
	fmt.Println(color.FgLightGreen.Render(string(k2ASCII)))

	msg := fmt.Sprintf("Initializing K2 KB System, version: %v build: %v PID: %v", version.GetVersion(), version.GetBuild(), os.Getppid())
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
	models.KBInit()
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
	go models.KBRun(&wg)
	go lib.Openbrowser("http://localhost" + os.Getenv("HTTPPORT"))
	wg.Wait()
}

func main() {
	//TEST

	//test.Test7()

	//CORE
	StartSystem()
}
