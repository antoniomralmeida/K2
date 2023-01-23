package locales

import (
	"github.com/antoniomralmeida/k2/initializers"
)

// Import these packages to trigger the init() function

// Locales is the list of locales's tags and names
// Please check if the language is supported in https://github.com/tebeka/snowball,
// if it is please add the correct language name.

const Locale_default = "en"

func InitStem() {
	for key := range initializers.Locales {
		st, err := initializers.NewStem(key)
		if err != nil {
			initializers.Log(err, initializers.Error)
		} else {
			l := initializers.Locales[key]
			l.Stemmer = st
			initializers.Locales[key] = l
		}
	}
}

// GetNameByTag returns the name of the given locale's tag
func GetNameByTag(tag string) string {
	return initializers.Locales[tag].Description
}

// GetTagByName returns the tag of the given locale's name
func GetTagByName(name string) string {

	for key, value := range initializers.Locales {
		if value.Description == name {
			return key
		}
	}
	return Locale_default
}

func GetLocaleByName(name string) initializers.Locale {
	for key, value := range initializers.Locales {
		if value.Description == name {
			return initializers.Locales[key]
		}
	}
	return initializers.Locales[Locale_default]
}

// Exists checks if the given tag exists in the list of locales
func Exists(tag string) bool {
	_, ok := initializers.Locales[tag]
	return ok
}
