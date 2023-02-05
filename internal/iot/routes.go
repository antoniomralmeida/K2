package iot

import "github.com/antoniomralmeida/k2/vendor/github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	app.Get("/api/v1/gates/:id", ReadGate)
	app.Post("/api/v1/gates/:id", WriteGate)
	app.Post("/api/v1/gates", NewGate)
	app.Get("/api/v1/gates", func(c *fiber.Ctx) error { return nil })
}
