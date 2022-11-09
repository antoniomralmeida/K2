package apikernel

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/kb"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func GetDataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	return c.JSON(objs)
}

func SetAttributeValue(c *fiber.Ctx) error {
	type Data struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var data Data
	err := json.Unmarshal(c.Body(), &data)
	lib.Log(err)
	a := kbbase.FindAttributeObjectByName(data.Name)
	if a != nil {
		a.SetValue(data.Value, kb.User, 100)
	} else {
		lib.Log(errors.New("Object not found! " + data.Name))
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.SendStatus(fiber.StatusOK)
}
