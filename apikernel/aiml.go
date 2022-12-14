package apikernel

import (
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/gofiber/fiber/v2"
)

var uid string

func GetChats(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	msg := c.Query("msg")
	uid := c.Query(fiber.HeaderXRequestID)
	//TODO: usar jwt como identificador de contexto
	//TODO: expandir aiml.xml

	if uid == "" || msg == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	resp := initializers.GetResponse(uid, msg)
	return c.Send([]byte(resp))
}
