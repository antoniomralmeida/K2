package web

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetDataInput(c *fiber.Ctx) error {
	objs := kbbase.GetDataInput()
	return c.JSON(objs)
}

func PostDataInput(c *fiber.Ctx) error {
	fields := c.FormValue("fileds")
	fmt.Println(fields)

	//TODO: Prsing objects and parsist

	return Home(c)
}
