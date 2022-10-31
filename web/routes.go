package web

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Static("/css", "./web/pub/css")
	app.Static("/img", "./web/pub/img")
	app.Static("/js", "./web/pub/js")
	app.Static("/vendor", "./web/pub/vendor")
	app.Static("/scss", "./web/pub/scss")

	app.Get("/", Home)
	app.Get("/api-datainput/", GetDataInput)
	app.Post("/api-datainput/", PostDataInput)
}
