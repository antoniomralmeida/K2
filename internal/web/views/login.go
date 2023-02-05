package views

import (
	"html/template"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/gofiber/fiber/v2"
)

func LoginView(c *fiber.Ctx) error {
	t, err := template.ParseFiles(T["login"].minify)
	if err != nil {
		inits.Log(err, inits.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.SendFile(T["404"].minify)
	}
	model := template.Must(t, nil)
	inits.Log(model.Execute(c, context.Ctxweb), inits.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	//c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	return c.SendStatus(fiber.StatusOK)
}
