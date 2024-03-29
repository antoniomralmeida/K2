package inits

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/antoniomralmeida/golibretranslate"
	"github.com/antoniomralmeida/k2/internal/lib"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const (
	I18n_accessforbidden     = "I18n_accessforbidden"
	I18n_alert               = "I18n_alert"
	I18n_already             = "I18n_already"
	I18n_alreadyregistered   = "I18n_alreadyregistered"
	I18n_avatar              = "I18n_avatar"
	I18n_badrequest          = "I18n_badrequest"
	I18n_cancel              = "I18n_cancel"
	I18n_dateinput           = "I18n_dateinput"
	I18n_forgot              = "I18n_forgot"
	I18n_halt                = "I18n_halt"
	I18n_internalservererror = "I18n_internalservererror"
	I18n_invalidcredentials  = "I18n_invalidcredentials"
	I18n_invalidimage        = "I18n_invalidimage"
	I18n_register            = "I18n_register"
	I18n_remember            = "I18n_remember"
	I18n_resume              = "I18n_resume"
	I18n_send                = "I18n_send"
	I18n_title               = "I18n_title"
	I18n_wellcome            = "I18n_wellcome"
	I18n_wellcome2           = "I18n_wellcome2"
	I18n_workspace           = "I18n_workspace"
	I18n_lgpd_cookie_consent = "I18n_lgpd_cookie_consent"
	I18n_lgpd_cookie_title   = "I18n_lgpd_cookie_title"
)

type Locale struct {
	Description string
	Voice       string
	Stemmer     *Stem
}

var (
	DefaultLocale = language.English.String()
	Locales       map[string]Locale
	bundle        *i18n.Bundle
	I18n_en       map[string]string
	TTSPathDir    = lib.GetWorkDir() + "/web/tts"
	TTSURL        = "tts/"
)

func I18nTranslate(original *map[string]string, locale string) (map[string]string, error) {
	translated := make(map[string]string)
	for key := range *original {
		trans, err := golibretranslate.Translate((*original)[key], DefaultLocale, locale)
		if err == nil {
			translated[key] = trans
		} else {
			return translated, err
		}
	}
	return translated, nil
}

func inLocalesConfig(locale string) bool {
	if os.Getenv("LOCALES") == "" {
		return true
	} else {
		locales := strings.Split(os.Getenv("LOCALES"), "|")
		for _, v := range locales {
			if v == locale {
				return true
			}
		}
	}
	return false
}

func GetSupportedLocales() (ret string) {

	ret = ""
	for key, value := range Locales {
		ret = ret + value.Description + "[" + key + "] "
	}
	return
}
func NewSupportedLanguage(locale string) {
	if inLocalesConfig(locale) || locale == language.English.String() || locale == language.Portuguese.String() {
		toEN := display.English.Languages()
		tag := language.MustParse(locale)
		Locales[locale] = Locale{Description: toEN.Name(tag)}
	}
}

func InitLangs() {
	Locales = make(map[string]Locale)

	NewSupportedLanguage(DefaultLocale)
	//
	NewSupportedLanguage(language.Portuguese.String())
	NewSupportedLanguage(language.Spanish.String())
	NewSupportedLanguage(language.German.String())
	NewSupportedLanguage(language.Hindi.String())
	NewSupportedLanguage(language.Arabic.String())
	NewSupportedLanguage(language.Bengali.String())
	NewSupportedLanguage(language.Russian.String())
	NewSupportedLanguage(language.Japanese.String())
	NewSupportedLanguage(language.French.String())
	NewSupportedLanguage(language.Italian.String())
	NewSupportedLanguage(language.Chinese.String())
	NewSupportedLanguage(language.Greek.String())
	NewSupportedLanguage(language.Dutch.String())

	I18n_en = make(map[string]string)

	I18n_en[I18n_accessforbidden] = "Access forbidden!"
	I18n_en[I18n_alert] = "Alerts"
	I18n_en[I18n_already] = "Already have an account? Login!"
	I18n_en[I18n_alreadyregistered] = "Already registered!"
	I18n_en[I18n_avatar] = "Avatar"
	I18n_en[I18n_badrequest] = "Bad Request!"
	I18n_en[I18n_cancel] = "Cancel"
	I18n_en[I18n_dateinput] = "Data Input"
	I18n_en[I18n_forgot] = "Forgot Password?"
	I18n_en[I18n_halt] = "KnowledgeBase was halt! "
	I18n_en[I18n_internalservererror] = "Internal Server Error!"
	I18n_en[I18n_invalidcredentials] = "Invalid credentials!"
	I18n_en[I18n_invalidimage] = "Invalid or empty image!"
	I18n_en[I18n_register] = "Please register to access K2!"
	I18n_en[I18n_remember] = "Remember me"
	I18n_en[I18n_resume] = "KnowledgeBase was resume! "
	I18n_en[I18n_send] = "Send"
	I18n_en[I18n_title] = "K2 KnowledgeBase System "
	I18n_en[I18n_wellcome] = "Wellcome to K2!"
	I18n_en[I18n_wellcome2] = "What are we going to do today?"
	I18n_en[I18n_workspace] = "Workspace"
	I18n_en[I18n_lgpd_cookie_consent] = "We use cookies to optimize your experience using the K2 system, only to store non-personal data on the application status, application usage preferences and information security in the browser. By clicking \"Accept All\", you agree to the use of cookies on this system."
	I18n_en[I18n_lgpd_cookie_title] = "We value your privacy"

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for code := range Locales {
		tomFile := TomlFile(code)
		_, err := os.Stat(tomFile)
		if code == DefaultLocale || os.IsNotExist(err) {
			if code == DefaultLocale {
				js, err := json.MarshalIndent(I18n_en, "", "	")
				Log(err, Error)
				f, err := os.Create(tomFile)
				Log(err, Error)
				f.WriteString(string(js))
				f.Close()
			} else {
				i18n, err := I18nTranslate(&I18n_en, code)
				Log(err, Fatal)
				f, err := os.Create(tomFile)
				Log(err, Error)
				js, err := json.MarshalIndent(i18n, "", "	")
				Log(err, Fatal)
				f.WriteString(string(js))
				f.Close()
			}
		}

		bundle.MustLoadMessageFile(tomFile)
	}
}

func TranslateTag(tag string, langs string) string {
	localizer := i18n.NewLocalizer(bundle, langs)
	msg, err := localizer.LocalizeMessage(&i18n.Message{ID: tag})
	Log(err, Error)
	return msg
}

func TomlFile(code string) string {
	wd := lib.GetWorkDir()
	path := wd + "/data/locales/"
	if ok, _ := lib.Exists(path); !ok {
		err := os.MkdirAll(path, os.ModePerm)
		Log(err, Fatal)
	}
	return path + "i18n." + code + ".json"
}

type Stem struct {
	stams map[string][]string
}

func (s *Stem) Stem(stem string) []string {
	return s.stams[stem]
}

func NewStem(locale string) (*Stem, error) {
	stam := new(Stem)
	stam.stams = make(map[string][]string)

	err, filename := lib.DownloadFile("https://raw.githubusercontent.com/michmech/lemmatization-lists/master/lemmatization-"+locale+".txt", lib.GetWorkDir()+"/data/locales/"+locale+"/")
	if err != nil {
		return nil, err
	}
	readFile, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) == 2 {
			stam.stams[fields[0]] = append(stam.stams[fields[0]], fields[1])
		}
	}

	readFile.Close()
	return stam, nil
}
