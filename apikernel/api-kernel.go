package apikernel

import (
	"log"
	"os"
	"sync"

	"github.com/antoniomralmeida/k2/kb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var kbbase *kb.KnowledgeBased

func Run(wg *sync.WaitGroup, kb *kb.KnowledgeBased) {
	defer wg.Done()
	kbbase = kb

	app := fiber.New(fiber.Config{AppName: "K2 System API-KERNEL v1.0.1",
		DisableStartupMessage: false,
		Prefork:               false})

	wd, _ := os.Getwd()
	f := wd + os.Getenv("HTTPAPILOG")
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v %v", err, f)
	}
	defer file.Close()

	app.Use(logger.New(logger.Config{Output: file,
		TimeFormat: "02/01/2006 15:04:05",
		Format:     "${time} [${ip}:${port}] ${status} ${latency} ${method} ${path} \n"}))
	app.Use(requestid.New())
	// Provide a custom compression level
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	Routes(app)

	app.Listen(os.Getenv("APIPORT"))
	wg.Done()
}
