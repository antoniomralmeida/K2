package context

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifyAvatar(c *fiber.Ctx) bool {
	if len(c.Query("avatar")) == 0 && len(Ctxweb.Avatar) > 0 && len(c.Cookies("avatar")) == 0 {
		cookie := fiber.Cookie{
			Name:     "avatar",
			Value:    Ctxweb.Avatar,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		}
		c.Cookie(&cookie)
		url := c.BaseURL()
		if c.Query("lang") != "" {
			url += "?lang=" + c.Query("lang")
		}
		c.Redirect(url)
		return true
	}
	return false
}
