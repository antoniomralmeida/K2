package web

import (
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/itmisx/i18n"
)

var supported_langs = []string{"en", "pt"}

func InitLangs() {
	var langPack1 = map[string]map[interface{}]interface{}{
		"en": {
			"title":     "K2 System KnowledgeBase",
			"datainput": "Data Input",
			"user":      "User",
			"workspace": "Workspace",
			"alerts":    "Alerts",
			"wellcome":  "Welcome Back!",
		},
		"pt": {
			"title":     "K2 System KnowledgeBase",
			"datainput": "Entrada de Dados",
			"user":      "Usuário",
			"workspace": "Área de trabalho",
			"alerts":    "Alertas",
			"wellcome":  "Bem Vindo de volta!",
		},
	}

	i18n.LoadLangPack(langPack1)
}

func Translate(term string, AcceptLanguage string) string {
	start := 0
	end := 0
	for i := 0; i < len(AcceptLanguage); i++ {
		if AcceptLanguage[i] == ',' {
			start = i + 1
		}
		if AcceptLanguage[i] == ';' {
			end = i
			break
		}
	}
	lang := AcceptLanguage[start:end]
	ok := false
	for _, l := range supported_langs {
		if lang == l {
			ok = true
			break
		}
	}
	if !ok {
		initializers.Log("Accept-Language "+AcceptLanguage+" not supported!", initializers.Info)
		lang = "en"
	}
	return i18n.T(lang, term)
}
