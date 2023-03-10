package views

import (
	"html/template"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/gofiber/fiber/v2"
)

func HomeView(c *fiber.Ctx) error {
	t, err := template.ParseFiles(T["home"].original)
	if err != nil {
		inits.Log(err, inits.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return ErrorView(c, err)
	}
	model := template.Must(t, nil)
	inits.Log(model.Execute(c, context.Ctxweb), inits.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)
}
