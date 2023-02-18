package main

import (
	"fmt"

	"os"
	"sync"

	"github.com/gookit/color"

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
	ctx, spanbase := inits.Begin("k2-init", nil)
	_, span := inits.Begin("inits.ConnectDB()", ctx)
	inits.ConnectDB()
	span.End()
	_, span = inits.Begin("models.KBInit()", ctx)
	models.InitKB()
	span.End()
	_, span = inits.Begin("inits.InitMQTT()", ctx)
	inits.InitMQTT()
	span.End()
	_, span = inits.Begin("models.InitSecurity()", ctx)
	models.InitSecurity()
	span.End()

	spanbase.End()
}

func StartSystem() {

	// CORE
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(5)
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
