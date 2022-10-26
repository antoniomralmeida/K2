package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Divs struct {
	User       string
	Title      string
	Workspacee string
}

func Run() {

	app := fiber.New()
	app.Use(logger.New())
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

	app.Listen(":3000")
}
