package main

import (
	"fmt"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/k2olivia/util"
	"github.com/antoniomralmeida/k2/k2web/web"
	"github.com/antoniomralmeida/k2/version"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/color"
)

func init() {

	k2webASCII := string(util.ReadFile("./config/k2web.txt"))
	fmt.Println(color.FgLightGreen.Render(k2webASCII))

	msg := fmt.Sprintf("Initializing Web Server K2 system version: %v build: %v", version.GetVersion(), version.GetBuild())
	initializers.InitEnvVars()
	initializers.InitLangs()
	initializers.LogInit("k2weblog")
	if !fiber.IsChild() {
		fmt.Println(msg)
		fmt.Println("Supported Languages: " + initializers.GetSupportedLocales())
		initializers.Log(msg, initializers.Info)
	}

}

func main() {
	web.Run()
}
