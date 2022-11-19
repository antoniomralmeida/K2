package apikernel

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/api/getlistdatainput", GetDataInput)
	app.Post("/api/postdatainput", PostDataInput)
	app.Get("/api/getlistworkspaces", GetWorkspaces)
}
