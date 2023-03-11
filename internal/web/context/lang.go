package context

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/pkg/htgotts"

	"github.com/gofiber/fiber/v2"
)

func getCurrentLang(c *fiber.Ctx) string {
	lang := c.Query("lang")
	accept := c.GetReqHeaders()["Accept-Language"]
	LangQ := ParseAcceptLanguage(lang, accept)
	for _, l := range LangQ {
		if _, ok := inits.Locales[l.Lang]; ok {
			return l.Lang
		}
	}
	return inits.DefaultLocale
}

func SetContextInfo(c *fiber.Ctx, templateMsg string) error {
	if Ctxweb.I18n == nil {
		Ctxweb.I18n = make(map[string]string)
	}
	var currentLang string
	currentLang = getCurrentLang(c)
	for key := range inits.I18n_en {
		Ctxweb.I18n[key] = inits.TranslateTag(key, currentLang)
	}

	if Ctxweb.Locales == nil {
		Ctxweb.Locales = make(map[string]string)
		for key, value := range inits.Locales {
			Ctxweb.Locales[key] = value.Description
		}
	}
	if templateMsg != "" {
		model, err := template.ParseFiles(templateMsg)
		if err != nil {
			return err
		}
		model = template.Must(model, nil)
		var tpl bytes.Buffer
		err = model.Execute(&tpl, Ctxweb)
		if err != nil {
			return err
		}
		text := tpl.String()
		file, err := htgotts.TTSToFile(text, currentLang, inits.TTSPathDir)
		if err != nil {
			return err
		}
		Ctxweb.WellcomeMsg = inits.TTSURL + filepath.Base(file)
	}
	keys := lib.DecodeToken(c.Cookies("jwt"))
	Ctxweb.User = fmt.Sprintf("%s", keys["name"])
	Ctxweb.UserId = fmt.Sprintf("%s", keys["user_id"])
	return nil
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
