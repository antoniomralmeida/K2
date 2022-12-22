package web

import (
	"github.com/BurntSushi/toml"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/models"
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

type i18n_Message struct {
	i18n_title     string
	i18n_wellcome  string
	i18n_datainput string
	i18n_wellcome2 string
	i18n_workspace string
	i18n_alert     string
	i18n_register  string
	i18n_already   string
	i18n_forgot    string
}

var i18n_en = i18n_Message{"K2 System KnowledgeBase",
	"Wellcome to K2!",
	"Data Input",
	"What are we going to do today?",
	"Workspace",
	"Alerts",
	"Please register to access K2!",
	"Already have an account? Login!",
	"Forgot Password?"}

func InitLangs() {

	models.Languages = append(models.Languages, models.Language{Code: "en", Description: "English", SpeechSynthesisId: 1})
	models.Languages = append(models.Languages, models.Language{Code: "pt", Description: "Português Brasileiro", SpeechSynthesisId: 0})
	//models.Languages = append(models.Languages, models.Language{Code: "es", Description: "Espanhol", SpeechSynthesisId: 262})
	//models.Languages = append(models.Languages, models.Language{Code: "de", Description: "Germany", SpeechSynthesisId: 143})

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	bundle.MustAddMessages(language.English, i18n.MustNewMessage(i18n_en))

	for _, l := range models.Languages {
		if l.Code != "en" {
			bundle.MustLoadMessageFile(TomlFile(l.Code))
		}
	}
}

func Translate(id string, c *fiber.Ctx) string {
	lang := c.Query("lang")
	accept := c.GetReqHeaders()["Accept-Language"]
	localizer := i18n.NewLocalizer(bundle, lang, accept)
	msg, err := localizer.LocalizeMessage(&i18n.Message{ID: id})
	initializers.Log(err, initializers.Error)
	return msg
}

func TomlFile(code string) string {
	return "./config/i18n." + code + ".json"
}
