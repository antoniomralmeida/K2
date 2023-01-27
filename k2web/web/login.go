package web

import (
	"html/template"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/models"
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
		initializers.Log(err, initializers.Error)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.SendFile(T["404"].minify)
	}
	model := template.Must(t, nil)
	initializers.Log(model.Execute(c, ctxweb), initializers.Error)
	c.Response().Header.Add("Content-Type", "text/html")
	//c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	return c.SendStatus(fiber.StatusOK)
}

func PostLogin(c *fiber.Ctx) error {
	req := models.LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		msg := translateID("i18n_badrequest", c) + ":" + err.Error()
		initializers.Log(msg, initializers.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}
	if req.Email == "" || req.Password == "" {
		msg := translateID("i18n_invalidcredentials", c)
		initializers.Log(msg, initializers.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	user := models.KBUser{}
	err := user.FindOne(bson.D{{Key: "email", Value: req.Email}})
	if err != mongo.ErrNoDocuments {
		initializers.Log(err, initializers.Error)
	}
	if user.Email == "" {
		msg := translateID("i18n_invalidcredentials", c)
		initializers.Log(msg, initializers.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password)); err != nil {
		msg := translateID("i18n_invalidcredentials", c)
		initializers.Log(msg, initializers.Info)
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	if user.Profile == models.Empty {
		msg := translateID("i18n_accessforbidden", c)
		initializers.Log(msg, initializers.Info)
		return fiber.NewError(fiber.StatusForbidden, msg)
	}

	token, _, err := lib.CreateJWTToken(user.ID, user.Name)
	if err != nil {
		msg := translateID("i18n_internalservererror", c)
		initializers.Log(msg, initializers.Info)
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
