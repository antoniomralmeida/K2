package apikernel

import (
	"github.com/antoniomralmeida/k2/cmd/k2/kb"
	"github.com/gofiber/fiber/v2"
)

func GetWorkspaces(c *fiber.Ctx) error {
	objs := kb.GKB.GetWorkspaces()
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	c.Response().Header.SetContentType("application/json")
	return c.Send([]byte(objs))
}
