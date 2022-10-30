package web

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	lang := c.GetReqHeaders()["Accept-Language"]
	ctxweb.Title = Translate("title", lang)
	ctxweb.DataInput = Translate("datainput", lang)
	ctxweb.User = "Manoel Ribeiro"
	ctxweb.Workspace = Translate("workspace", lang)
	ctxweb.Alerts = Translate("alerts", lang)

	model := template.Must(template.ParseFiles("./web/assets/template.html"))
	model.Execute(c, ctxweb)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)
}
