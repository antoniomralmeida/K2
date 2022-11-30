package iot

import (
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadGate(c *fiber.Ctx) error {
	// Read the param noteId
	id := c.Params("id")
	return c.SendStatus(fiber.StatusOK)
}

func WriteGate(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.SendStatus(fiber.StatusOK)
}

func NewGate(c *fiber.Ctx) error {
	gate := Gate{}
	c.BodyParser(gate)
	gate.Id = primitive.NewObjectID()
	ctx, collection := initializers.GetCollection("K2Gate")
	_, err := collection.InsertOne(ctx, gate)
	if err != nil {
		initializers.Log(err, initializers.Error)
		return c.SendStatus(fiber.StatusBadGateway)
	} else {
		return c.JSON(gate)
	}
}
