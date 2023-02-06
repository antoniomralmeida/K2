package analysis

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	"github.com/antoniomralmeida/golibretranslate"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/olivia/modules"
	"github.com/antoniomralmeida/k2/internal/olivia/util"
)

var intents = map[string][]Intent{}

// Intent is a way to group sentences that mean the same thing and link them with a tag which
// represents what they mean, some responses that the bot can reply and a context
type Intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Context   string   `json:"context"`
}

// Document is any sentence from the intents' patterns linked with its tag
type Document struct {
	Sentence Sentence
	Tag      string
}

// CacheIntents set the given intents to the global variable intents
func CacheIntents(locale string, _intents []Intent) {
	intents[locale] = _intents
}

// GetIntents returns the cached intents
func GetIntents(locale string) []Intent {
	return intents[locale]
}

func intentsFile(locale string) string {

	return inits.GetHomeDir() + "/data/locales/" + locale + "/intents.json"
}

func translateIntents(_intents *[]Intent, locale string) (err error) {
	for i := range *_intents {
		var trans string
		trans, err = golibretranslate.Translate((*_intents)[i].Tag, "en", locale)
		if err == nil {
			(*_intents)[i].Tag = trans
		} else {
			return
		}
		trans, err = golibretranslate.Translate((*_intents)[i].Context, "en", locale)
		if err == nil {
			(*_intents)[i].Context = trans
		} else {
			return
		}
		for j := range (*_intents)[i].Patterns {
			trans, err = golibretranslate.Translate((*_intents)[i].Patterns[j], "en", locale)
			if err == nil {
				(*_intents)[i].Patterns[j] = trans
			} else {
				return
			}
		}
		for j := range (*_intents)[i].Responses {
			trans, err = golibretranslate.Translate((*_intents)[i].Responses[j], "en", locale)
			if err == nil {
				(*_intents)[i].Responses[j] = trans
			} else {
				return
			}
		}
	}
	return
}

// SerializeIntents returns a list of intents retrieved from the given intents file
func SerializeIntents(locale string) (_intents []Intent) {
	intents := intentsFile(locale)
	if ok, _ := lib.Exists(intents); !ok {
		intents_tmp := []Intent{}
		intentsFile := intentsFile(inits.DefaultLocale)
		bytes, _ := ioutil.ReadFile(intentsFile)
		err := json.Unmarshal(bytes, &intents_tmp)
		inits.Log(err, inits.Fatal)
		err = translateIntents(&intents_tmp, locale)
		inits.Log(err, inits.Fatal)
		js, err := json.Marshal(intents_tmp)
		inits.Log(err, inits.Error)
		f, err := os.Create(intents)
		inits.Log(err, inits.Error)
		f.WriteString(string(js))
		f.Close()
	}
	txt, _ := ioutil.ReadFile(intents)
	err := json.Unmarshal(txt, &_intents)
	inits.Log(err, inits.Fatal)

	CacheIntents(locale, _intents)

	return _intents
}

// SerializeModulesIntents retrieves all the registered modules and returns an array of Intents
func SerializeModulesIntents(locale string) []Intent {
	registeredModules := modules.GetModules(locale)
	intents := make([]Intent, len(registeredModules))

	for k, module := range registeredModules {
		intents[k] = Intent{
			Tag:       module.Tag,
			Patterns:  module.Patterns,
			Responses: module.Responses,
			Context:   "",
		}
	}

	return intents
}

// GetIntentByTag returns an intent found by given tag and locale
func GetIntentByTag(tag, locale string) Intent {
	for _, intent := range GetIntents(locale) {
		if tag != intent.Tag {
			continue
		}

		return intent
	}

	return Intent{}
}

// Organize intents with an array of all words, an array with a representative word of each tag
// and an array of Documents which contains a word list associated with a tag
func Organize(locale string) (words, classes []string, documents []Document) {
	// Append the modules intents to the intents from res/datasets/intents.json
	intents := append(
		SerializeIntents(locale),
		SerializeModulesIntents(locale)...,
	)

	for _, intent := range intents {
		for _, pattern := range intent.Patterns {
			// Tokenize the pattern's sentence
			patternSentence := Sentence{locale, pattern}
			patternSentence.arrange()

			// Add each word to response
			for _, word := range patternSentence.stem() {

				if !util.Contains(words, word) {
					words = append(words, word)
				}
			}

			// Add a new document
			documents = append(documents, Document{
				patternSentence,
				intent.Tag,
			})
		}

		// Add the intent tag to classes
		classes = append(classes, intent.Tag)
	}

	sort.Strings(words)
	sort.Strings(classes)

	return words, classes, documents
}
