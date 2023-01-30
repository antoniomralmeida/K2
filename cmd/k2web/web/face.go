package web

import (
	"html/template"

	"github.com/antoniomralmeida/k2/inits"
	"github.com/gofiber/fiber/v2"
)

func GetFace(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		return c.Redirect(c.OriginalURL() + "?avatar=" + ctxweb.Avatar)
	}
	//Context
	SetContextInfo(c)

	//Render
	t, err := template.ParseFiles(T["face"].minify)
	if err != nil {
		inits.Log(err, inits.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.SendFile(T["404"].minify)
	}
	model := template.Must(t, nil)
	inits.Log(model.Execute(c, ctxweb), inits.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)
}
