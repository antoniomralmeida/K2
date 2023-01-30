package iot

import (
	"fmt"

	"github.com/antoniomralmeida/k2/inits"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
)

func ReadGate(c *fiber.Ctx) error {
	// Read the param noteId
	id := c.Params("id")
	fmt.Println(id)
	return c.SendStatus(fiber.StatusOK)
}

func WriteGate(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println(id)
	return c.SendStatus(fiber.StatusOK)
}

func NewGate(c *fiber.Ctx) error {
	gate := new(Gate)
	c.BodyParser(gate)
	err := mgm.Coll(gate).Create(gate)
	if err != nil {
		inits.Log(err, inits.Error)
		return c.SendStatus(fiber.StatusBadGateway)
	} else {
		return c.JSON(gate)
	}
}
