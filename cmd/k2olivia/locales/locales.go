package locales

import "github.com/antoniomralmeida/k2/inits"

// Import these packages to trigger the init() function

// Locales is the list of locales's tags and names
// Please check if the language is supported in https://github.com/tebeka/snowball,
// if it is please add the correct language name.

func InitStem() {
	for key := range inits.Locales {
		st, err := inits.NewStem(key)
		if err != nil {
			inits.Log(err, inits.Error)
		} else {
			l := inits.Locales[key]
			l.Stemmer = st
			inits.Locales[key] = l
		}
	}
}

// GetNameByTag returns the name of the given locale's tag
func GetNameByTag(tag string) string {
	return inits.Locales[tag].Description
}

// GetTagByName returns the tag of the given locale's name
func GetTagByName(name string) string {

	for key, value := range inits.Locales {
		if value.Description == name {
			return key
		}
	}
	return inits.DefaultLocale
}

func GetLocaleByName(name string) inits.Locale {
	for key, value := range inits.Locales {
		if value.Description == name {
			return inits.Locales[key]
		}
	}
	return inits.Locales[inits.DefaultLocale]
}

// Exists checks if the given tag exists in the list of locales
func Exists(tag string) bool {
	_, ok := inits.Locales[tag]
	return ok
}
