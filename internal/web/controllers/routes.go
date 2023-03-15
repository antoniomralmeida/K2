package controllers

import (
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

func Routes(app *fiber.App) {
	app.Static("/css", "./web/css")
	app.Static("/img", "./web/img")
	app.Static("/upload", "./web/upload")
	app.Static("/js", "./web/js")
	app.Static("/vendor", "./web/vendor")
	app.Static("/scss", "./web/scss")
	app.Static("/tts", "./web/tts")
	app.Static("/audio", "./web/audio")

	// Allow cors for cookie
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Use(csrf.New(csrf.Config{
		KeyLookup:  "form:csrf",
		CookieName: "csrf_",
		ContextKey: "csrf",
	}))

	app.Get("/", Splash)
	app.Get("/home", Home)
	app.Get("/login", LoginForm)
	app.Post("/login", PostLogin)
	app.Post("/logout", Logout)
	app.Get("/face", GetFace)
	app.Get("/signup", SignUpForm)
	app.Post("/signup", PostSignUp)

	api := app.Group("/api/v1")
	api.Get("/attributes", GetAttributes)
	api.Post("/attributes", PostAttributes)
	api.Get("/chats", GetChats)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return Error(c, lib.PageNotFoundError)
	})
}
