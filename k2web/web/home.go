package web

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	if lib.ValidateToken(c.Cookies("jwt")) != nil {
		c.SendStatus(fiber.StatusForbidden)
		return c.Redirect("/login")
	}
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		return c.Redirect(c.BaseURL() + "?avatar=" + ctxweb.Avatar)
	}
	//Context
	TranslateAll(c)

	call := ctxweb.ApiKernel + "/workspaces"
	resp, err := http.Get(call)
	if err != nil {
		initializers.Log(err, initializers.Error)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		initializers.Log(err, initializers.Error)
		err = json.Unmarshal(body, &ctxweb.Workspaces)
		initializers.Log(err, initializers.Error)
	}
	//Render
	model := template.Must(template.ParseFiles(T["home"].minify))
	model.Execute(c, ctxweb)
	c.Response().Header.Add("Content-Type", "text/html")

	return c.SendStatus(fiber.StatusOK)
}
