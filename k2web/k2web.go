package main

import (
	"fmt"
	"log"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/k2web/web"
	"github.com/antoniomralmeida/k2/version"
	"github.com/gofiber/fiber/v2"
	"github.com/subosito/gotenv"
)

func init() {
	msg := fmt.Sprintf("Initializing Web Server K2 system version: %v build: %v", version.Version, version.Build)

	if err := gotenv.Load("./.env"); err != nil {
		log.Fatal(err)
	}
	initializers.LogInit("k2weblog")
	if !fiber.IsChild() {
		fmt.Println(msg)
		initializers.Log(msg, initializers.Info)
	}
}
func main() {
	web.Run()
}
