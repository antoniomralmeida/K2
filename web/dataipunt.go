package web

import (
	"strings"

	"github.com/antoniomralmeida/k2/kb"
	"github.com/gofiber/fiber/v2"
)

func GetDataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	return c.JSON(objs)
}

func PostDataInput(c *fiber.Ctx) error {
	fields := c.FormValue("fields")
	fs := strings.Split(fields, "|")
	for i := range fs {
		if len(fs[i]) > 0 {
			a := kbbase.FindAttributeObjectByName(fs[i])
			if a != nil {
				a.SetValue(c.FormValue(fs[i]), kb.KBSource(kb.User), 100)
			}
		}
	}
	c.Append("Location", "/")
	return Home(c)
}
