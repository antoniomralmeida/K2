package main

import (
	"fmt"
	"os"

	"github.com/antoniomralmeida/k2/cmd/iotmidware/iot"
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
)

func init() {
	msg := fmt.Sprintf("Initializing K2 IoT Midware: %v build: %v PID: %v", lib.GetVersion(), lib.GetBuild(), os.Getppid())
	fmt.Println(msg)
	inits.InitEnvVars()
	inits.LogInit("k2iot")
	inits.ConnectDB()
}

func main() {
	iot.Run()
}
