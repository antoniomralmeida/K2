package controllers

import (
	"github.com/antoniomralmeida/k2/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetWorkspaces(c *fiber.Ctx) error {
	objs := models.KBWorkspacesJson()
	c.Response().Header.SetContentType("application/json")
	return c.Send([]byte(objs))
}
