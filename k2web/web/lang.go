package web

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

var I18n_ID = []string{
	"i18n_title",
	"i18n_wellcome",
	"i18n_dateinput",
	"i18n_wellcome2",
	"i18n_workspace",
	"i18n_alert",
	"i18n_register",
	"i18n_already",
	"i18n_forgot",
	"i18n_send",
	"i18n_cancel"}

type I18n_Messages struct {
	I18n_title     string `json:"i18n_title"`
	I18n_wellcome  string `json:"i18n_wellcome"`
	I18n_dateinput string `json:"i18n_dateinput"`
	I18n_wellcome2 string `json:"i18n_wellcome2"`
	I18n_workspace string `json:"i18n_workspace"`
	I18n_alert     string `json:"i18n_alert"`
	I18n_register  string `json:"i18n_register"`
	I18n_already   string `json:"i18n_already"`
	I18n_forgot    string `json:"i18n_forgot"`
	I18n_send      string `json:"i18n_forgot"`
	I18n_cancel    string `json:"i18n_forgot"`
}

var i18n_en = I18n_Messages{
	"K2 System KnowledgeBase",
	"Wellcome to K2!",
	"Data Input",
	"What are we going to do today?",
	"Workspace",
	"Alerts",
	"Please register to access K2!",
	"Already have an account? Login!",
	"Forgot Password?",
	"Send",
	"Cancel"}

type Language struct {
	Code              string
	Description       string
	SpeechSynthesisId int
}

var Languages []Language

func InitLangs() {

	Languages = append(Languages, Language{Code: "en", Description: "English", SpeechSynthesisId: 1})
	Languages = append(Languages, Language{Code: "pt", Description: "Português Brasileiro", SpeechSynthesisId: 0})
	Languages = append(Languages, Language{Code: "es", Description: "Espanhol", SpeechSynthesisId: 262})
	Languages = append(Languages, Language{Code: "de", Description: "Germany", SpeechSynthesisId: 143})

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, l := range Languages {
		_, err := os.Stat(TomlFile(l.Code))
		if l.Code == "en" || os.IsNotExist(err) {
			f, err := os.Create(TomlFile(l.Code))
			initializers.Log(err, initializers.Error)
			js, err := json.Marshal(i18n_en)
			initializers.Log(err, initializers.Error)
			f.WriteString(string(js))
			f.Close()
		}
		bundle.MustLoadMessageFile(TomlFile(l.Code))
	}
}

func SetContextInfo(c *fiber.Ctx) {
	if ctxweb.I18n == nil {
		ctxweb.I18n = make(map[string]string)
	}
	for _, id := range I18n_ID {
		ctxweb.I18n[id] = translateID(id, c)
	}
	lang := c.Query("lang")
	accept := c.GetReqHeaders()["Accept-Language"]
	LangQ := ParseAcceptLanguage(lang, accept)
	for _, l := range Languages {
		for _, l2 := range LangQ {
			if l2.Lang == l.Code {
				ctxweb.SpeechSynthesisId = l.SpeechSynthesisId
				break
			}
		}
	}
	ctxweb.JwtToken = c.Cookies("jwt")
	keys := lib.DecodeToken(ctxweb.JwtToken)
	ctxweb.User = fmt.Sprintf("%s", keys["name"])
	ctxweb.UserId = fmt.Sprintf("%s", keys["user_id"])
}

func translateID(id string, c *fiber.Ctx) string {
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

type LangQ struct {
	Lang string
	Q    float64
}

func ParseAcceptLanguage(lang, acptLang string) []LangQ {
	var lqs []LangQ
	if lang != "" {
		lqs = append(lqs, LangQ{lang, 1})
	}
	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		langQ := strings.Split(trimedLangQStr, ";")
		if len(langQ) == 1 {
			lq := LangQ{langQ[0], 1}
			lqs = append(lqs, lq)
		} else {
			qp := strings.Split(langQ[1], "=")
			q, err := strconv.ParseFloat(qp[1], 64)
			if err != nil {
				initializers.Log(err, initializers.Fatal)

			}
			lq := LangQ{langQ[0], q}
			lqs = append(lqs, lq)
		}
	}
	return lqs
}
