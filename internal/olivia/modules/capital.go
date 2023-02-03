package modules

import (
	"fmt"
	"regexp"

	"github.com/antoniomralmeida/k2/internal/olivia/language"
	"github.com/antoniomralmeida/k2/internal/olivia/util"
)

var (
// CapitalTag is the intent tag for its module
// CapitalTag = "capital"
// ArticleCountries is the map of functions to find the article in front of a country
// in different languages
// ArticleCountries = map[string]func(string) string{}
)

// CapitalReplacer replaces the pattern contained inside the response by the capital of the country
// specified in the message.
// See modules/modules.go#Module.Replacer() for more details.
func CapitalReplacer(locale, entry, response, _ string) (string, string) {
	country := language.FindCountry(locale, entry)

	// If there isn't a country respond with a message from res/datasets/messages.json
	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, util.GetMessage(locale, responseTag)
	}
	countryName := ArticleCountries(locale, country.Name[locale])
	/*
		articleFunction, exists := ArticleCountries[locale]
		countryName := country.Name[locale]
		if exists {
			countryName = articleFunction(countryName)
		}
	*/
	return CapitalTag, fmt.Sprintf(response, countryName, country.Capital)
}

func ArticleCountries(locale, name string) string {
	arts, exists := articles[locale]
	if exists {
		for _, a := range arts {
			match, _ := regexp.MatchString(a.Regexp, name)
			if match {
				return a.Article + " " + name
			}
		}
	}
	return name
}
