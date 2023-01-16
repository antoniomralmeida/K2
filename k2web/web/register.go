package web

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignupForm(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		url := c.OriginalURL()
		sep := "?"
		if strings.Contains(url, sep) {
			sep = "&"
		}
		return c.Redirect(url + sep + "avatar=" + ctxweb.Avatar)
	}
	//Context
	SetContextInfo(c)
	//Render
	t, err := template.ParseFiles(T["sigup"].minify)
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

func PostSignup(c *fiber.Ctx) error {
	req := new(models.SigupRequest)
	fmt.Println(string(c.Request().Header.ContentType()))
	fmt.Println(c.FormValue("name"))
	fmt.Println(c.FormValue("email"))
	fmt.Println(c.FormValue("password"))
	fmt.Println(c.FormValue("password2"))

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, translateID("i18n_badrequest", c)+":"+err.Error())
	}
	faceimage := "faceimage"
	file, err := c.FormFile(faceimage)
	if err != nil {
		initializers.Log(err, initializers.Error)
		return fiber.NewError(fiber.StatusBadRequest, translateID("i18n_invalidimage", c))
	}
	if file != nil {
		filePath := lib.TempFileName("k2-upload", filepath.Ext(file.Filename))
		err = c.SaveFile(file, filePath)
		if err != nil {
			initializers.Log(err, initializers.Error)
			return fiber.NewError(fiber.StatusInternalServerError, translateID("i18n_internalservererror", c))
		}
		faceimage = filePath
	} else {
		faceimage = ""
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, translateID("i18n_invalidcredentials", c))
	}

	if req.Password2 != req.Password {
		return fiber.NewError(fiber.StatusBadRequest, translateID("i18n_invalidcredentials", c))
	}

	user := models.KBUser{}
	err = user.FindOne(bson.D{{Key: "email", Value: req.Email}})
	if err != mongo.ErrNoDocuments && err != nil {
		initializers.Log(err, initializers.Error)
		return fiber.NewError(fiber.StatusInternalServerError, translateID("i18n_internalservererror", c))
	}

	if user.Email != "" {
		return fiber.NewError(fiber.StatusBadRequest, translateID("i18n_alreadyregistered", c))
	}

	err = models.NewUser(req.Name, req.Email, req.Password, faceimage)
	if err != nil {
		initializers.Log(err, initializers.Error)
		return fiber.NewError(fiber.StatusInternalServerError, translateID("i18n_internalservererror", c)+":"+err.Error())
	}
	if faceimage != "" {
		os.Remove(faceimage)
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.Redirect("/login")
}
