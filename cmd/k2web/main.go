package main

import (
	"fmt"
	"os"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/web"
	"github.com/gofiber/fiber/v2"

	"github.com/gookit/color"
)

func init() {

	k2webASCII, _ := os.ReadFile("./configs/k2web.txt")
	fmt.Println(color.FgLightGreen.Render(string(k2webASCII)))

	msg := fmt.Sprintf("Initializing Web Server from K2 KB System,  version: %v build: %v", lib.GetVersion(), lib.GetBuild())
	inits.InitEnvVars()
	inits.InitLangs()
	inits.LogInit("k2weblog")
	if !fiber.IsChild() {
		fmt.Println(msg)
		fmt.Println("Supported Languages: " + inits.GetSupportedLocales())
		inits.Log(msg, inits.Info)
	}

}

func main() {
	web.Run()
}
