package web

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Divs struct {
	User       string
	Title      string
	Workspacee string
}

var runnig = false

func Run() {

	if runnig {
		return
	}

	app := fiber.New()

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

	app.Get("/", func(c *fiber.Ctx) error {
		d := Divs{Title: "teste"}

		return c.Render("./web/assets/gomodel.html", d)
	})
	app.Post("/api-*", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("*"))
	})
	runnig = true
	app.Listen(":3000")
}
