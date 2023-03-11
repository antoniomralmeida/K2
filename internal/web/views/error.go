package views

import (
	"github.com/antoniomralmeida/k2/internal/web/context"

	"github.com/gofiber/fiber/v2"
)

func ErrorView(c *fiber.Ctx, err error) error {
	context.Ctxweb.Error = err.Error()
	return c.Render(T["error"].original, context.Ctxweb)
}
