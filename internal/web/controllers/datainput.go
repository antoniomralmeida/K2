package controllers

import (
	"net/url"

	"github.com/antoniomralmeida/k2/internal/fuzzy"
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetAttributes(c *fiber.Ctx) error {
	objs := models.KBGetDataInput()
	return c.JSON(objs)
}

func PostAttributes(c *fiber.Ctx) error {

	data, err := url.ParseQuery(string(c.Body()))
	if inits.Log(err, inits.Error) != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	for key := range data {
		a := models.KBFindAttributeObjectByName(key)
		if a != nil {
			a.SetValue(data[key][0], models.FromUser, fuzzy.TrustUser)
		} else {
			inits.Log("Object not found! "+key, inits.Error)
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
