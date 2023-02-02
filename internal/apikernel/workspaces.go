package apikernel

import (
	"github.com/antoniomralmeida/k2/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetWorkspaces(c *fiber.Ctx) error {
	objs := models.KBGetWorkspaces()
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	c.Response().Header.SetContentType("application/json")
	return c.Send([]byte(objs))
}
