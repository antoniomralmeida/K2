package web

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/antoniomralmeida/golibretranslate"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

var i18n_en map[string]string

func InitLangs() {
	i18n_en = make(map[string]string)

	i18n_en[inits.I18n_accessforbidden] = "Access forbidden!"
	i18n_en[inits.I18n_alert] = "Alerts"
	i18n_en[inits.I18n_already] = "Already have an account? Login!"
	i18n_en[inits.I18n_alreadyregistered] = "Already registered!"
	i18n_en[inits.I18n_avatar] = "Avatar"
	i18n_en[inits.I18n_badrequest] = "Bad Request!"
	i18n_en[inits.I18n_cancel] = "Cancel"
	i18n_en[inits.I18n_dateinput] = "Data Input"
	i18n_en[inits.I18n_forgot] = "Forgot Password?"
	i18n_en[inits.I18n_halt] = "KnowledgeBase was halt! "
	i18n_en[inits.I18n_internalservererror] = "Internal Server Error!"
	i18n_en[inits.I18n_invalidcredentials] = "Invalid credentials!"
	i18n_en[inits.I18n_invalidimage] = "Invalid or empty image!"
	i18n_en[inits.I18n_register] = "Please register to access K2!"
	i18n_en[inits.I18n_remember] = "Remember me"
	i18n_en[inits.I18n_resume] = "KnowledgeBase was resume! "
	i18n_en[inits.I18n_send] = "Send"
	i18n_en[inits.I18n_title] = "K2 KnowledgeBase System "
	i18n_en[inits.I18n_wellcome] = "Wellcome to K2!"
	i18n_en[inits.I18n_wellcome2] = "What are we going to do today?"
	i18n_en[inits.I18n_workspace] = "Workspace"

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for code := range inits.Locales {
		tomFile := TomlFile(code)
		_, err := os.Stat(tomFile)
		if code == inits.DefaultLocale || os.IsNotExist(err) {
			if code == inits.DefaultLocale {

				js, err := json.Marshal(i18n_en)
				inits.Log(err, inits.Error)
				f, err := os.Create(tomFile)
				inits.Log(err, inits.Error)
				f.WriteString(string(js))
				f.Close()
			} else {
				i18n, err := I18nTranslate(&i18n_en, code)
				inits.Log(err, inits.Fatal)
				f, err := os.Create(tomFile)
				inits.Log(err, inits.Error)
				js, err := json.Marshal(i18n)
				inits.Log(err, inits.Fatal)
				f.WriteString(string(js))
				f.Close()
			}
		}
		bundle.MustLoadMessageFile(tomFile)
	}
}

func I18nTranslate(orignal *map[string]string, locale string) (map[string]string, error) {
	translated := make(map[string]string)
	for key := range *orignal {
		trans, err := golibretranslate.Translate((*orignal)[key], inits.DefaultLocale, locale)
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
		if l2, ok := inits.Locales[l.Lang]; ok {
			ctxweb.SpeechSynthesisId = l2.SpeechSynthesisId
			break
		}
	}
	ctxweb.Locales = make(map[string]string)
	for key, value := range inits.Locales {
		ctxweb.Locales[key] = value.Description
	}
	//ctxweb.JwtToken = c.Cookies("jwt")
	keys := lib.DecodeToken(c.Cookies("jwt"))
	ctxweb.User = fmt.Sprintf("%s", keys["name"])
	ctxweb.UserId = fmt.Sprintf("%s", keys["user_id"])
}

func translateID(id string, c *fiber.Ctx) string {
	lang := c.Query("lang")
	accept := c.GetReqHeaders()["Accept-Language"]
	localizer := i18n.NewLocalizer(bundle, lang, accept)
	msg, err := localizer.LocalizeMessage(&i18n.Message{ID: id})
	inits.Log(err, inits.Error)
	return msg
}

func TomlFile(code string) string {
	wd := inits.GetHomeDir()
	path := wd + "/web/res/locale/"
	if ok, _ := lib.Exists(path); !ok {
		err := os.MkdirAll(path, os.ModePerm)
		inits.Log(err, inits.Fatal)
	}
	return path + "i18n." + code + ".json"
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
				inits.Log(err, inits.Fatal)

			}
			lq := LangQ{langQ[0], q}
			lqs = append(lqs, lq)
		}
	}
	return lqs
}
