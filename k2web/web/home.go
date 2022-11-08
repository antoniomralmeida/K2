package web

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	//Context
	lang := c.GetReqHeaders()["Accept-Language"]
	ctxweb.Title = Translate("title", lang)
	ctxweb.DataInput = Translate("datainput", lang)
	ctxweb.User = "Manoel Ribeiro"
	ctxweb.Workspace = Translate("workspace", lang)
	ctxweb.Alerts = Translate("alerts", lang)
	//Render

	model := template.Must(template.ParseFiles(T["home"].original))
	model.Execute(c, ctxweb)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)

	//return c.Render(T["home"].original, ctxweb)
}
