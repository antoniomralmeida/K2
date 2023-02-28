package context

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifyCookies(c *fiber.Ctx) bool {
	if len(c.Cookies("avatar")) == 0 || len(c.Cookies("lang")) == 0 {
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
		lang := "pt"
		if len(c.Query("lang")) > 0 {
			lang = c.Query("lang")
		}
		cookie2 := fiber.Cookie{
			Name:     "lang",
			Value:    lang,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: false,
		}
		c.Cookie(&cookie2)
		c.Redirect(c.BaseURL())
	}

	return false
}
