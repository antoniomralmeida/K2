package web

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Static("/css", GetK2Path()+"/k2web/web/pub/css")
	app.Static("/img", GetK2Path()+"/k2web/web/pub/img")
	app.Static("/js", GetK2Path()+"/k2web/web/pub/js")
	app.Static("/vendor", GetK2Path()+"/k2web/web/pub/vendor")
	app.Static("/scss", GetK2Path()+"/k2web/web/pub/scss")

	app.Get("/", Home)
	app.Get("/api-datainput/", GetDataInput)
	app.Post("/api-datainput/", PostDataInput)
}
