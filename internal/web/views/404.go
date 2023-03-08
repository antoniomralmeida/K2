package views

import (
	"github.com/antoniomralmeida/k2/internal/web/context"

	"github.com/gofiber/fiber/v2"
)

func NotFoundView(c *fiber.Ctx) error {
	c.Render(T["404"].original, context.Ctxweb)
	return c.SendStatus(fiber.StatusNotFound)
}
