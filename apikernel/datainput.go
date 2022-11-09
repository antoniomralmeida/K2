package apikernel

import (
	"net/url"

	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func GetDataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	return c.JSON(objs)
}

func PostDataInput(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	data, err := url.ParseQuery(string(c.Body()))
	lib.Log(err)
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	for key := range data {
		a := kbbase.FindAttributeObjectByName(key)
		if a != nil {
			a.SetValue(data[key][0], kb.User, 100)
		} else {
			lib.Log(errors.New("Object not found! " + key))
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
