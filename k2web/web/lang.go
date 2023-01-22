package web

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/antoniomralmeida/golibretranslate"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

type Language struct {
	Description       string
	SpeechSynthesisId int
}

var languages map[string]Language
var i18n_en map[string]string

func InitLangs() {
	languages = make(map[string]Language)
	languages["en"] = Language{Description: "English", SpeechSynthesisId: 1}
	languages["pt"] = Language{Description: "Portuguese(BR)", SpeechSynthesisId: 0}
	languages["es"] = Language{Description: "Spanish", SpeechSynthesisId: 262}
	languages["de"] = Language{Description: "German", SpeechSynthesisId: 143}
	languages["hi"] = Language{Description: "Hindi", SpeechSynthesisId: 154}
	languages["ar"] = Language{Description: "Arabic", SpeechSynthesisId: 12}
	languages["bn"] = Language{Description: "Bengali", SpeechSynthesisId: 48}
	languages["ru"] = Language{Description: "Russian", SpeechSynthesisId: 213}
	languages["ja"] = Language{Description: "Japanese", SpeechSynthesisId: 167}
	languages["fr"] = Language{Description: "French", SpeechSynthesisId: 133}
	languages["it"] = Language{Description: "Italian", SpeechSynthesisId: 164}

	i18n_en = make(map[string]string)
	i18n_en["i18n_title"] = "K2 System KnowledgeBase"
	i18n_en["i18n_wellcome"] = "Wellcome to K2!"
	i18n_en["i18n_dateinput"] = "Data Input"
	i18n_en["i18n_wellcome2"] = "What are we going to do today?"
	i18n_en["i18n_workspace"] = "Workspace"
	i18n_en["i18n_alert"] = "Alerts"
	i18n_en["i18n_register"] = "Please register to access K2!"
	i18n_en["i18n_already"] = "Already have an account? Login!"
	i18n_en["i18n_forgot"] = "Forgot Password?"
	i18n_en["i18n_send"] = "Send"
	i18n_en["i18n_cancel"] = "Cancel"
	i18n_en["i18n_remember"] = "Remember me"
	i18n_en["i18n_badrequest"] = "Bad Request"
	i18n_en["i18n_invalidimage"] = "Invalid or empty image"
	i18n_en["i18n_internalservererror"] = "InternalServerError"
	i18n_en["i18n_invalidcredentials"] = "Invalid credentials"
	i18n_en["i18n_alreadyregistered"] = "Already registered"
	i18n_en["i18n_accessforbidden"] = "Access forbidden"

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for code := range languages {
		_, err := os.Stat(TomlFile(code))
		if code == "en" || os.IsNotExist(err) {
			f, err := os.Create(TomlFile(code))
			initializers.Log(err, initializers.Error)
			if code == "en" {
				js, err := json.Marshal(i18n_en)
				initializers.Log(err, initializers.Error)
				f.WriteString(string(js))
			} else {
				i18n, err := I18nTranslate(&i18n_en, code)
				initializers.Log(err, initializers.Fatal)
				js, err := json.Marshal(i18n)
				initializers.Log(err, initializers.Fatal)
				f.WriteString(string(js))
			}
			f.Close()
		}
		bundle.MustLoadMessageFile(TomlFile(code))
	}
}

func I18nTranslate(orignal *map[string]string, locale string) (map[string]string, error) {
	translated := make(map[string]string)
	for key, _ := range *orignal {
		fmt.Println(key)
		trans, err := golibretranslate.Translate((*orignal)[key], "en", locale)
		if err == nil {
			translated[key] = trans
		} else {
			return translated, err
		}
	}
	return translated, nil
}

func SetContextInfo(c *fiber.Ctx) {
	if ctxweb.I18n == nil {
		ctxweb.I18n = make(map[string]string)
	}
	for key := range i18n_en {
		ctxweb.I18n[key] = translateID(key, c)
	}
	lang := c.Query("lang")
	accept := c.GetReqHeaders()["Accept-Language"]
	LangQ := ParseAcceptLanguage(lang, accept)
	for _, l := range LangQ {
		if l2, ok := languages[l.Lang]; ok {
			ctxweb.SpeechSynthesisId = l2.SpeechSynthesisId
			break
		}
	}
	ctxweb.Locales = make(map[string]string)
	for key, value := range languages {
		ctxweb.Locales[key] = value.Description
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
	path := "./k2web/pub/res/locale/" + code + "/"
	if ok, _ := lib.Exists(path); !ok {
		err := os.Mkdir(path, os.FileMode(0777))
		initializers.Log(err, initializers.Fatal)
	}
	return path + "i18n.json"
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
