package modules

import (
	"encoding/json"
	"reflect"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/olivia/util"
)

const (
	JokesTag          = "JokesTag"
	AdvicesTag        = "AdvicesTag"
	AreaTag           = "AreaTag"
	CapitalTag        = "CapitalTag"
	CurrencyTag       = "CurrencyTag"
	GenresTag         = "GenresTag"
	MoviesTag         = "MoviesTag"
	MoviesAlreadyTag  = "MoviesAlreadyTag"
	MoviesDataTag     = "MoviesDataTag"
	NameGetterTag     = "NameGetterTag"
	NameSetterTag     = "NameSetterTag"
	RandomTag         = "RandomTag"
	ReminderSetterTag = "ReminderSetterTag"
	ReminderGetterTag = "ReminderGetterTag"
	SpotifySetterTag  = "SpotifySetterTag"
	SpotifyPlayerTag  = "SpotifyPlayerTag"
)

// Module is a structure for dynamic intents with a Tag, some Patterns and Responses and
// a Replacer function to execute the dynamic changes.
type Module struct {
	Tag             string
	Patterns        []string
	Responses       []string
	Replacer        string
	Context         string
	ReflectReplacer reflect.Value
	//func(string, string, string, string) (string, string)
}

type Article struct {
	Regexp  string
	Article string
}

type My struct{}

var (
	modules  = map[string][]Module{}
	articles = map[string][]Article{}
)

func modulesFile(locale string) string {
	return inits.GetHomeDir() + "/data/locales/" + locale + "/modules.json"
}

func articleFile(locale string) string {
	return inits.GetHomeDir() + "/data/locales/" + locale + "/article.json"
}

func init() {

	for locale := range inits.Locales {

		mfile := modulesFile(locale)
		if ok, _ := lib.Exists(mfile); ok {
			var mods []Module
			err := json.Unmarshal(util.ReadFile(mfile), &mods)
			inits.Log(err, inits.Fatal)

			for i := range mods {
				m := My{}
				mods[i].ReflectReplacer = reflect.ValueOf(m).MethodByName(mods[i].Replacer)
			}

			RegisterModules(locale, mods)
		}
		afile := articleFile(locale)
		if ok, _ := lib.Exists(afile); ok {
			var atrs []Article
			err := json.Unmarshal(util.ReadFile(afile), &atrs)
			inits.Log(err, inits.Fatal)
			RegisterArticles(locale, atrs)
		}
	}

}

// RegisterModule registers a module into the map
func RegisterModule(locale string, module Module) {
	modules[locale] = append(modules[locale], module)
}

// RegisterModules registers an array of modules into the map
func RegisterModules(locale string, _modules []Module) {
	modules[locale] = append(modules[locale], _modules...)
}

func RegisterArticles(locale string, _atrs []Article) {
	articles[locale] = append(articles[locale], _atrs...)
}

// GetModules returns all the registered modules
func GetModules(locale string) []Module {
	return modules[locale]
}

// GetModuleByTag returns a module found by the given tag and locale
func GetModuleByTag(tag, locale string) Module {
	for _, module := range modules[locale] {
		if tag != module.Tag {
			continue
		}

		return module
	}

	return Module{}
}

// ReplaceContent apply the Replacer of the matching module to the response and returns it
func ReplaceContent(locale, tag, entry, response, token string) (string, string) {
	for _, module := range modules[locale] {
		if module.Tag != tag {
			continue
		}
		if module.ReflectReplacer.IsValid() {
			params := []reflect.Value{reflect.ValueOf(locale), reflect.ValueOf(entry), reflect.ValueOf(response), reflect.ValueOf(token)}
			ret := module.ReflectReplacer.Call(params)
			return ret[0].String(), ret[1].String()
		} else {
			return "", ""
		}
	}

	return tag, response
}
