package apikernel

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/api/v1/attributes", GetDataInput)
	app.Post("/api/v1/attributes", PostDataInput)
	app.Get("/api/v1/workspaces", GetWorkspaces)
}
