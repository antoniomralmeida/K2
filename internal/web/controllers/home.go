package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/antoniomralmeida/k2/internal/web/views"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	if lib.ValidateToken(c.Cookies("jwt")) != nil {
		c.SendStatus(fiber.StatusForbidden)
		url := c.BaseURL() + "/login?lang=" + c.Query("lang") + "&avatar=" + context.Ctxweb.Avatar
		return c.Redirect(url)
	}
	if len(c.Query("avatar")) == 0 && len(context.Ctxweb.Avatar) > 0 {
		url := c.BaseURL() + "?lang=" + c.Query("lang") + "&avatar=" + context.Ctxweb.Avatar
		return c.Redirect(url)
	}
	//Context
	context.SetContextInfo(c)
	call := context.Ctxweb.ApiKernel + "/workspaces"
	resp, err := http.Get(call)
	if err != nil {
		inits.Log(err, inits.Error)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		inits.Log(err, inits.Error)
		err = json.Unmarshal(body, &context.Ctxweb.Workspaces)
		inits.Log(err, inits.Error)
	}

	return views.HomeView(c)
}
