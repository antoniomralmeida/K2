package apikernel

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Get("/attributes", GetAttributes)
	api.Post("/attributes", PostAttributes)
	api.Get("/workspaces", GetWorkspaces)
	api.Get("/chats", GetChats)
}
