package controllers

import (
	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/gofiber/fiber/v2"
)

func GetChats(c *fiber.Ctx) error {

	msg := c.Query("msg")
	uid := c.Query("jwt")
	lang := c.Query("lang")

	if msg == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	resp := inits.GetResponse(lang, uid, msg)
	return c.Send([]byte(resp))
}
