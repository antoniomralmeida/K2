package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/models"
	"github.com/antoniomralmeida/k2/services"

	"github.com/antoniomralmeida/k2/version"
)

func init() {
	msg := fmt.Sprintf("Initializing K2 System version: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
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
	initializers.InitOlivia()
	span.End()
	spanbase.End()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(5)
	go apikernel.Run(&wg)
	go kb.Run(&wg)
	go services.Run(&wg)
	go lib.Openbrowser("http://localhost" + os.Getenv("HTTPPORT"))
	wg.Wait()
}

func main() {
	//TEST

	//test.Test7()

	//CORE
	StartSystem()
}
