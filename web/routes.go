package web

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	app.Get("/", Home)
	app.Post("/api-datainput*", DataInput)
}
