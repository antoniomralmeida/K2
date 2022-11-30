package main

import (
	"fmt"
	"os"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/iotmidware/iot"
	"github.com/antoniomralmeida/k2/version"
)

func init() {
	msg := fmt.Sprintf("Initializing K2 IoT Midware: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
	fmt.Println(msg)
	initializers.InitEnvVars()
	initializers.LogInit("k2iot")
	initializers.ConnectDB()
}

func main() {
	iot.Run()
}
