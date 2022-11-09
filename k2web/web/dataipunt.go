package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
)

func GetDataInput(c *fiber.Ctx) error {
	if apikernel != "" {
		callapi := apikernel + "/getlistdatainput"
		resp, err := http.Get(callapi)
		lib.LogFatal(err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		lib.LogFatal(err)
		return c.Send(body)
	}
	return fiber.ErrBadGateway
}

func PostDataInput(c *fiber.Ctx) error {
	//application/x-www-form-urlencoded
	data, err := url.ParseQuery(string(c.Body()))
	lib.LogFatal(err)
	for key := range data {
		callapi := apikernel + "/setattributevalue"
		body, err := json.Marshal(map[string]string{"name": key, "value": data[key][0]})
		fmt.Println(callapi, body, err)
		lib.LogFatal(err)
		http.Post(callapi, "application/json", bytes.NewBuffer(body))
	}
	c.Append("Location", "/")
	return Home(c)
}
