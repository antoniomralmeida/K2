package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Routes(app *fiber.App) {
	app.Static("/css", "./k2web/pub/css")
	app.Static("/img", "./k2web/pub/img")
	app.Static("/upload", "./k2web/pub/upload")
	app.Static("/js", "./k2web/pub/js")
	app.Static("/vendor", "./k2web/pub/vendor")
	app.Static("/scss", "./k2web/pub/scss")
	app.Static("/sounds", "./k2web/pub/sounds")

	// Allow cors for cookie
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/", Home)

	app.Get("/login", LoginForm)
	app.Post("/login", PostLogin)
	app.Post("/logout", Logout)
	app.Get("/face", GetFace)
	app.Get("/signup", SignupForm)
	app.Post("/signup", PostSignup)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		c.Render(T["404"].original, nil)
		return c.SendStatus(fiber.StatusNotFound)
	})
}
