package web

import (
	"html/template"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/gofiber/fiber/v2"
)

func GetFace(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		return c.Redirect(c.OriginalURL() + "?avatar=" + ctxweb.Avatar)
	}
	//Context
	lang := c.GetReqHeaders()["Accept-Language"]
	ctxweb.Title = Translate("title", lang)

	model := template.Must(template.ParseFiles(T["face"].original))
	initializers.Log(model.Execute(c, ctxweb), initializers.Error)
	c.Response().Header.Add("Content-Type", "text/html")

	return c.SendStatus(fiber.StatusOK)
}
