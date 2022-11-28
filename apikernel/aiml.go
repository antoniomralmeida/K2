package apikernel

import (
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/gofiber/fiber/v2"
)

var uid string

func GetChats(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	c.Request().Header.VisitAll(func(key, value []byte) {
		if string(key) == "X-Request-Id" {
			uid = string(value)
		}
	})
	line := c.Query("line")
	return c.Send([]byte(initializers.GetResponse(uid, line)))
}
