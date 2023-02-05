package apikernel

import (
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/gofiber/fiber/v2"
)

func GetChats(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")

	msg := c.Query("msg")
	uid := c.Query("jwt")
	lang := c.Query("lang")

	if msg == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	resp := inits.GetResponse(lang, uid, msg)
	return c.Send([]byte(resp))
}
