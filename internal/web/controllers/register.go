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

func SignUpForm(c *fiber.Ctx) error {
	if context.VerifyCookies(c) {
		return nil
	}
	//Context
	err := context.SetContextInfo(c, views.T["register_wellcome"].FullPath)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.ErrInternalServerError
	}
	return views.RegisterView(c)
}

func PostSignUp(c *fiber.Ctx) error {
	req := new(models.SigupRequest)
	context.SetContextInfo(c, "")
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, context.Ctxweb.I18n[inits.I18n_badrequest]+":"+err.Error())
	}
	faceimage := "faceimage"
	file, err := c.FormFile(faceimage)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusBadRequest, context.Ctxweb.I18n[inits.I18n_invalidimage])
	}
	if file != nil {
		filePath := lib.TempFileName("k2-upload", filepath.Ext(file.Filename))
		err = c.SaveFile(file, filePath)
		if err != nil {
			inits.Log(err, inits.Error)
			return fiber.NewError(fiber.StatusInternalServerError, context.Ctxweb.I18n[inits.I18n_internalservererror])
		}
		faceimage = filePath
	} else {
		faceimage = ""
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, context.Ctxweb.I18n[inits.I18n_invalidcredentials])
	}

	if req.Password2 != req.Password {
		return fiber.NewError(fiber.StatusBadRequest, context.Ctxweb.I18n[inits.I18n_invalidcredentials])
	}

	user := models.KBUser{}
	err = user.FindOne(bson.D{{Key: "email", Value: req.Email}})
	if err != mongo.ErrNoDocuments && err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusInternalServerError, context.Ctxweb.I18n[inits.I18n_internalservererror])
	}

	if user.Email != "" {
		return fiber.NewError(fiber.StatusBadRequest, context.Ctxweb.I18n[inits.I18n_alreadyregistered])
	}

	err = models.NewUser(req.Name, req.Email, req.Password, faceimage)
	if err != nil {
		inits.Log(err, inits.Error)
		return fiber.NewError(fiber.StatusInternalServerError, context.Ctxweb.I18n[inits.I18n_internalservererror]+":"+err.Error())
	}
	if faceimage != "" {
		os.Remove(faceimage)
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.Redirect("/login")
}
