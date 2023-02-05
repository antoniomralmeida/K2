package controllers

import (
	"os"
	"path/filepath"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/antoniomralmeida/k2/internal/web/views"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignupForm(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(context.Ctxweb.Avatar) > 0 {
		url := c.BaseURL() + "/signup?lang=" + c.Query("lang") + "&avatar=" + context.Ctxweb.Avatar
		return c.Redirect(url)
	}
	//Context
	context.SetContextInfo(c)
	return views.RegisterView(c)
}

func PostSignup(c *fiber.Ctx) error {
	req := new(models.SigupRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, context.TranslateTag(inits.I18n_badrequest, c)+":"+err.Error())
	}
	faceimage := "faceimage"
	file, err := c.FormFile(faceimage)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusBadRequest, context.TranslateTag(inits.I18n_invalidimage, c))
	}
	if file != nil {
		filePath := lib.TempFileName("k2-upload", filepath.Ext(file.Filename))
		err = c.SaveFile(file, filePath)
		if err != nil {
			inits.Log(err, inits.Error)
			return fiber.NewError(fiber.StatusInternalServerError, context.TranslateTag(inits.I18n_internalservererror, c))
		}
		faceimage = filePath
	} else {
		faceimage = ""
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, context.TranslateTag(inits.I18n_invalidcredentials, c))
	}

	if req.Password2 != req.Password {
		return fiber.NewError(fiber.StatusBadRequest, context.TranslateTag(inits.I18n_invalidcredentials, c))
	}

	user := models.KBUser{}
	err = user.FindOne(bson.D{{Key: "email", Value: req.Email}})
	if err != mongo.ErrNoDocuments && err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusInternalServerError, context.TranslateTag(inits.I18n_internalservererror, c))
	}

	if user.Email != "" {
		return fiber.NewError(fiber.StatusBadRequest, context.TranslateTag(inits.I18n_alreadyregistered, c))
	}

	err = models.NewUser(req.Name, req.Email, req.Password, faceimage)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusInternalServerError, context.TranslateTag(inits.I18n_internalservererror, c)+":"+err.Error())
	}
	if faceimage != "" {
		os.Remove(faceimage)
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.Redirect("/login")
}
