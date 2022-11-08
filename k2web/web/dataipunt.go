package web

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetDataInput(c *fiber.Ctx) error {

	if apikernel != "" {
		callapi := apikernel + "/getlistdatainput"
		api := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(api)

		req := api.Request()
		req.Header.SetMethod(fiber.MethodGet)
		req.SetRequestURI(callapi)
		//FIX: não está chamando a API
		if err := api.Parse(); err != nil {
			log.Println(err)
		} else {
			code, body, errs := api.Bytes()
			if errs != nil {
				return c.Send(body)
			} else {
				c.SendStatus(code)
			}
		}
	}
	return fiber.ErrBadGateway
}

func PostDataInput(c *fiber.Ctx) error {
	var data map[string]string
	c.BodyParser(&data)
	for key := range data {
		callapi := apikernel + "/setattributevalue"
		api := fiber.AcquireAgent()
		req := api.Request()
		req.Header.Add(key, data[key])
		req.Header.SetMethod("post")
		req.SetRequestURI(callapi)
		if err := api.Parse(); err != nil {
			log.Println(err)
		} else {
			if _, _, errs := api.Bytes(); errs == nil {
				log.Println(errs)
			}
		}
	}
	c.Append("Location", "/")
	return Home(c)
}
