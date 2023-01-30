package web

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignupForm(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		url := c.BaseURL() + "/signup?lang=" + c.Query("lang") + "&avatar=" + ctxweb.Avatar
		return c.Redirect(url)
	}
	//Context
	SetContextInfo(c)
	//Render
	t, err := template.ParseFiles(T["sigup"].minify)
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

func PostSignup(c *fiber.Ctx) error {
	req := new(models.SigupRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, translateID(inits.I18n_badrequest, c)+":"+err.Error())
	}
	faceimage := "faceimage"
	file, err := c.FormFile(faceimage)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusBadRequest, translateID(inits.I18n_invalidimage, c))
	}
	if file != nil {
		filePath := lib.TempFileName("k2-upload", filepath.Ext(file.Filename))
		err = c.SaveFile(file, filePath)
		if err != nil {
			inits.Log(err, inits.Error)
			return fiber.NewError(fiber.StatusInternalServerError, translateID(inits.I18n_internalservererror, c))
		}
		faceimage = filePath
	} else {
		faceimage = ""
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, translateID(inits.I18n_invalidcredentials, c))
	}

	if req.Password2 != req.Password {
		return fiber.NewError(fiber.StatusBadRequest, translateID(inits.I18n_invalidcredentials, c))
	}

	user := models.KBUser{}
	err = user.FindOne(bson.D{{Key: "email", Value: req.Email}})
	if err != mongo.ErrNoDocuments && err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusInternalServerError, translateID(inits.I18n_internalservererror, c))
	}

	if user.Email != "" {
		return fiber.NewError(fiber.StatusBadRequest, translateID(inits.I18n_alreadyregistered, c))
	}

	err = models.NewUser(req.Name, req.Email, req.Password, faceimage)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusInternalServerError, translateID(inits.I18n_internalservererror, c)+":"+err.Error())
	}
	if faceimage != "" {
		os.Remove(faceimage)
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.Redirect("/login")
}
