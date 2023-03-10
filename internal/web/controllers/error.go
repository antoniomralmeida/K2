package controllers

import (
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/antoniomralmeida/k2/internal/web/views"
	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, err error) error {
	//Context
	context.SetContextInfo(c, "")
	return views.ErrorView(c, err)
}
