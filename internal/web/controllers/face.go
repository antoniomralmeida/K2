package controllers

import (
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/antoniomralmeida/k2/internal/web/views"
	"github.com/gofiber/fiber/v2"
)

func GetFace(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(context.Ctxweb.Avatar) > 0 {
		return c.Redirect(c.OriginalURL() + "?avatar=" + context.Ctxweb.Avatar)
	}
	//Context
	context.SetContextInfo(c, "")

	return views.FaceView(c)
}
