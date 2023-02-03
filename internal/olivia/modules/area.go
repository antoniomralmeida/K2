package modules

import (
	"fmt"

	"github.com/antoniomralmeida/k2/internal/olivia/language"
	"github.com/antoniomralmeida/k2/internal/olivia/util"
)

// AreaReplacer replaces the pattern contained inside the response by the area of the country
// specified in the message.
// See modules/modules.go#Module.Replacer() for more details.
func AreaReplacer(locale, entry, response, _ string) (string, string) {
	country := language.FindCountry(locale, entry)

	// If there isn't a country respond with a message from res/datasets/messages.json
	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	return AreaTag, fmt.Sprintf(response, ArticleCountries(locale, country.Name[locale]), country.Area)
}
