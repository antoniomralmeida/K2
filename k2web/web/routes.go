package web

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Static("/css", "./../k2web/pub/css")
	app.Static("/img", "./../k2web/pub/img")
	app.Static("/js", "./../k2web/pub/js")
	app.Static("/vendor", "./../k2web/pub/vendor")
	app.Static("/scss", "./../k2web/pub/scss")

	app.Get("/", Home)
}
