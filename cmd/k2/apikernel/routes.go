package apikernel

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/api/v1/attributes", GetAttributes)
	app.Post("/api/v1/attributes", PostAttributes)
	app.Get("/api/v1/workspaces", GetWorkspaces)
	app.Get("/api/v1/chats", GetChats)
}
