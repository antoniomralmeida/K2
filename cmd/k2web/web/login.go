package web

import (
	"html/template"
	"time"

	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginForm(c *fiber.Ctx) error {
	if len(c.Query("avatar")) == 0 && len(ctxweb.Avatar) > 0 {
		url := c.BaseURL() + "/login?lang=" + c.Query("lang") + "&avatar=" + ctxweb.Avatar
		return c.Redirect(url)
	}
	//Context
	SetContextInfo(c)

	//TODO: Incluir reconhecimento facil no login

	//Render
	t, err := template.ParseFiles(T["login"].minify)
	if err != nil {
		inits.Log(err, inits.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.SendFile(T["404"].minify)
	}
	model := template.Must(t, nil)
	inits.Log(model.Execute(c, ctxweb), inits.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	//c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	return c.SendStatus(fiber.StatusOK)
}

func PostLogin(c *fiber.Ctx) error {
	req := models.LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		msg := translateID(inits.I18n_badrequest, c) + ":" + err.Error()
		inits.Log(msg, inits.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	if req.Email == "" || req.Password == "" {
		msg := translateID(inits.I18n_invalidcredentials, c)
		inits.Log(msg, inits.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	user := models.KBUser{}
	err := user.FindOne(bson.D{{Key: "email", Value: req.Email}})
	if err != mongo.ErrNoDocuments {
		inits.Log(err, inits.Error)
	}
	if user.Email == "" {
		msg := translateID(inits.I18n_invalidcredentials, c)
		inits.Log(msg, inits.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password)); err != nil {
		msg := translateID(inits.I18n_invalidcredentials, c)
		inits.Log(msg, inits.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	if user.Profile == models.Empty {
		msg := translateID(inits.I18n_accessforbidden, c)
		inits.Log(msg, inits.Info)
		return fiber.NewError(fiber.StatusForbidden, msg)
	}

	token, _, err := lib.CreateJWTToken(user.ID, user.Name)
	if err != nil {
		msg := translateID(inits.I18n_internalservererror, c)
		inits.Log(msg, inits.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	ctxweb.User = user.Name
	// Create cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	c.SendStatus(fiber.StatusAccepted)
	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	// Remove cookie
	// -time.Hour = expires before one hour
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	ctxweb.User = ""
	c.Cookie(&cookie)
	c.SendStatus(fiber.StatusAccepted)
	return c.Redirect("/login")
}
