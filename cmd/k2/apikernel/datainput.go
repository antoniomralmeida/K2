package apikernel

import (
	"net/url"

	"github.com/antoniomralmeida/k2/cmd/k2/kb"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/gofiber/fiber/v2"
)

func GetAttributes(c *fiber.Ctx) error {
	objs := kb.GKB.GetDataInput()
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	return c.JSON(objs)
}

func PostAttributes(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	data, err := url.ParseQuery(string(c.Body()))
	if inits.Log(err, inits.Error) != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	for key := range data {
		a := kb.GKB.FindAttributeObjectByName(key)
		if a != nil {
			a.SetValue(data[key][0], kb.User, 100)
		} else {
			inits.Log("Object not found! "+key, inits.Error)
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
