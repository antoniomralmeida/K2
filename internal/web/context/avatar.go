package context

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifyCookies(c *fiber.Ctx) bool {
	if len(c.Cookies("avatar")) == 0 || len(c.Query("avatar")) > 0 {
		avatar := ""
		if len(c.Query("avatar")) > 0 {
			avatar = c.Query("avatar")
		} else {
			avatar = Ctxweb.Avatar
		}

		cookie := fiber.Cookie{
			Name:     "avatar",
			Value:    avatar,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: false,
		}
		c.Cookie(&cookie)

		c.Redirect(c.BaseURL())
		return true
	}

	return false
}
