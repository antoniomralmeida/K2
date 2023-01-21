package locales

import (
	"fmt"

	"github.com/antoniomralmeida/k2/k2olivia/util"
)

// Import these packages to trigger the init() function

// Locales is the list of locales's tags and names
// Please check if the language is supported in https://github.com/tebeka/snowball,
// if it is please add the correct language name.

var Locales = []Locale{
	{
		Tag:  "en",
		Name: "english",
	},
	{
		Tag:  "de",
		Name: "german",
	},
	{
		Tag:  "fr",
		Name: "french",
	},
	{
		Tag:  "es",
		Name: "spanish",
	},
	{
		Tag:  "hi",
		Name: "hindi",
	},
	{
		Tag:  "it",
		Name: "italian",
	},
	{
		Tag:  "ar",
		Name: "arabic",
	},
	{
		Tag:  "bn",
		Name: "bengali",
	},
	{
		Tag:  "ja",
		Name: "japanese",
	},
	{
		Tag:  "pt",
		Name: "portuguese",
	},
	{
		Tag:  "ru",
		Name: "russian",
	},
}

func init() {
	for i := range Locales {
		var err error
		Locales[i].Stemmer, err = util.NewStem(Locales[i].Tag)
		if err != nil {
			fmt.Println("Stemmer error", err)
			return
		}
	}

}

// A Locale is a registered locale in the file
type Locale struct {
	Tag     string
	Name    string
	Stemmer *util.Stem
}

// GetNameByTag returns the name of the given locale's tag
func GetNameByTag(tag string) string {
	for _, locale := range Locales {
		if locale.Tag != tag {
			continue
		}

		return locale.Name
	}

	return "English"
}

// GetTagByName returns the tag of the given locale's name
func GetTagByName(name string) string {
	for _, locale := range Locales {
		if locale.Name != name {
			continue
		}
		return locale.Tag
	}
	return "en"
}

func GetLocaleByName(name string) Locale {
	for _, locale := range Locales {
		if locale.Name != name {
			continue
		}
		return locale
	}
	return Locales[0]
}

// Exists checks if the given tag exists in the list of locales
func Exists(tag string) bool {
	for _, locale := range Locales {
		if locale.Tag == tag {
			return true
		}
	}

	return false
}
