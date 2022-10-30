package web

import (
	"github.com/gofiber/fiber/v2"
)

func DataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	return c.JSON(objs)
}
