package web

import (
	"log"
	"os"
	"sync"

	"github.com/antoniomralmeida/k2/kb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Context struct {
	User      string
	Title     string
	DataInput string
	Workspace string
	Alerts    string
}

var ctxweb = Context{}
var kbbase *kb.KnowledgeBase

func Run(wg *sync.WaitGroup, kb *kb.KnowledgeBase) {
	defer wg.Done()
	kbbase = kb

	Init()

	app := fiber.New(fiber.Config{AppName: "K2 System v1.0.1",
		DisableStartupMessage: true,
		Prefork:               true})

	file, err := os.OpenFile("./log/web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app.Use(logger.New(logger.Config{Output: file,
		TimeFormat: "02/01/2006 15:04:05",
		Format:     "${time} [${ip}:${port}] ${status} ${latency} ${method} ${path} \n"}))
	app.Use(requestid.New())
	app.Static("/css", "./web/assets/css")
	app.Static("/img", "./web/assets/img")
	app.Static("/js", "./web/assets/js")
	app.Static("/vendor", "./web/assets/vendor")

	Routes(app)

	app.Listen(":3000")
	wg.Done()
}
