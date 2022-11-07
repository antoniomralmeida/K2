package apikernel

import (
	"github.com/antoniomralmeida/k2/kb"
	"github.com/gofiber/fiber/v2"
)

func GetDataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	return c.JSON(objs)
}

func SetAttributeValue(c *fiber.Ctx) error {
	var data map[string]string
	c.BodyParser(&data)
	for key := range data {
		a := kbbase.FindAttributeObjectByName(key)
		if a != nil {
			a.SetValue(c.FormValue(data[key]), kb.KBSource(kb.User), 100)
		} else {
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
