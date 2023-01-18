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
		url := c.BaseURL() + "/login?lang=" + c.Query("lang") + "&avatar=" + ctxweb.Avatar
		return c.Redirect(url)
	}
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		url := c.BaseURL() + "?lang=" + c.Query("lang") + "&avatar=" + ctxweb.Avatar
		return c.Redirect(url)
	}
	//Context
	SetContextInfo(c)

	//Render
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
	t, err := template.ParseFiles(T["home"].original)
	if err != nil {
		initializers.Log(err, initializers.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.SendFile(T["404"].minify)
	}
	model := template.Must(t, nil)
	initializers.Log(model.Execute(c, ctxweb), initializers.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)
}
