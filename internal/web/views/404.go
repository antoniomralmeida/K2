package views

import "github.com/gofiber/fiber/v2"

func NotFoundView(c *fiber.Ctx) error {
	c.Render(T["404"].original, nil)
	return c.SendStatus(fiber.StatusNotFound)
}
