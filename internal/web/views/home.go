package views

import (
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/gofiber/fiber/v2"
)

func HomeView(c *fiber.Ctx) error {
	return c.Render(T["home"].Minify, context.Ctxweb)
}
