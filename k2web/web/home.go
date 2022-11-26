package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	//Context
	lang := c.GetReqHeaders()["Accept-Language"]
	ctxweb.Title = Translate("title", lang)
	ctxweb.DataInput = Translate("datainput", lang)
	ctxweb.User = "Manoel Ribeiro"
	ctxweb.Workspace = Translate("workspace", lang)
	ctxweb.Alerts = Translate("alerts", lang)

	call := ctxweb.ApiKernel + "/workspaces"
	resp, err := http.Get(call)
	initializers.Log(err, initializers.Error)
	body, err := ioutil.ReadAll(resp.Body)
	initializers.Log(err, initializers.Error)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &ctxweb.Workspaces)
	initializers.Log(err, initializers.Error)

	//Render
	model := template.Must(template.ParseFiles(T["home"].original))
	model.Execute(c, ctxweb)
	c.Response().Header.Add("Content-Type", "text/html")
	return c.SendStatus(fiber.StatusOK)
}
