package apikernel

import (
	"fmt"
	"net/url"

	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
)

func GetDataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	return c.JSON(objs)
}

func SetAttributeValue(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	data, err := url.ParseQuery(string(c.Body()))
	lib.LogFatal(err)
	for key := range data {
		fmt.Println(key)
		a := kbbase.FindAttributeObjectByName(key)
		if a != nil {
			a.SetValue(c.FormValue(data[key][0]), kb.KBSource(kb.User), 100)
		} else {
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
