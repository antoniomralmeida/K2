package web

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
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
		inits.Log(err, inits.Error)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		inits.Log(err, inits.Error)
		err = json.Unmarshal(body, &ctxweb.Workspaces)
		inits.Log(err, inits.Error)
	}
	//Render
	t, err := template.ParseFiles(T["home"].original)
	if err != nil {
		inits.Log(err, inits.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.SendFile(T["404"].minify)
	}
	model := template.Must(t, nil)
	inits.Log(model.Execute(c, ctxweb), inits.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)
}
