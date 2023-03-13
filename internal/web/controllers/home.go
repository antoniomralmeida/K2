package controllers

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/antoniomralmeida/k2/internal/web/views"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	if lib.ValidateToken(c.Cookies("jwt")) != nil {
		c.SendStatus(fiber.StatusForbidden)
		url := c.BaseURL() + "/login"
		return c.Redirect(url)
	}
	if context.VerifyCookies(c) {
		return nil
	}
	//Context
	err := context.SetContextInfo(c, lib.GetWorkDir()+"/web/home_wellcome.html")
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.ErrInternalServerError
	}

	works, err := models.KBWorkspacesJson()
	if err != nil {
		inits.Log(err, inits.Error)
		c.Status(fiber.StatusInternalServerError)
	}
	err = json.Unmarshal([]byte(works), &context.Ctxweb.Workspaces)
	if err != nil {
		inits.Log(err, inits.Error)
		c.Status(fiber.StatusInternalServerError)
	}
	return views.HomeView(c)
}
