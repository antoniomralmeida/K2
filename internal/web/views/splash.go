package views

import (
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/gofiber/fiber/v2"
)

func SplashView(c *fiber.Ctx) error {
	return c.Render(T["splash"].original, context.Ctxweb)
}
