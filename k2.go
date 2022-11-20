package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/antoniomralmeida/k2/apikernel"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/services"
	"github.com/antoniomralmeida/k2/telemetry"
	"github.com/antoniomralmeida/k2/version"
)

func init() {
	msg := fmt.Sprintf("Initializing K2 System version: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
	fmt.Println(msg)
	initializers.InitEnvVars()
	initializers.LogInit("k2log")
	initializers.Log(msg, initializers.Info)
	telemetry.Init()
	ctx, spanbase := telemetry.Begin(context.TODO(), "main-init")
	_, span := telemetry.Begin(ctx, "ConnectDB")
	initializers.ConnectDB()
	span.End()
	_, span = telemetry.Begin(ctx, "kb.Init")
	kb.Init()
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

	//test.Test8()

	//CORE
	StartSystem()
}
