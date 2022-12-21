package web

import (
	"github.com/BurntSushi/toml"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

const (
	i18n_title     = "title"
	i18n_wellcome  = "wellcome"
	i18n_wellcome2 = "wellcome2"
	i18n_dateinput = "datainput"
	i18n_workspace = "workspace"
	i18n_alert     = "alert"
	i18n_register  = "register"
)

func InitLangs() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("config/en.toml")
	bundle.MustLoadMessageFile("config/pt.toml")
}

func Translate(id string, c *fiber.Ctx) string {
	lang := c.Query("lang")
	accept := c.GetReqHeaders()["Accept-Language"]
	localizer := i18n.NewLocalizer(bundle, lang, accept)
	msg, err := localizer.LocalizeMessage(&i18n.Message{ID: id})
	initializers.Log(err, initializers.Error)
	return msg
}
