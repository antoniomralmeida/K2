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
	TranslateAll(c)

	model := template.Must(template.ParseFiles(T["face"].minify))
	initializers.Log(model.Execute(c, ctxweb), initializers.Error)
	c.Response().Header.Add("Content-Type", "text/html")

	return c.SendStatus(fiber.StatusOK)
}
