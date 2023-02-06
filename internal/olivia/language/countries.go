package language

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/antoniomralmeida/k2/internal/inits"
)

// Country is the serializer of the countries.json file in the res folder
type Country struct {
	Name     map[string]string `json:"name"`
	Capital  string            `json:"capital"`
	Code     string            `json:"code"`
	Area     float64           `json:"area"`
	Currency string            `json:"currency"`
}

var countries = SerializeCountries()

// SerializeCountries returns a list of countries, serialized from `res/datasets/countries.json`
func SerializeCountries() (countries []Country) {
	wd := inits.GetHomeDir()
	bytes, _ := ioutil.ReadFile(wd + "/data/datasets/countries.json")
	err := json.Unmarshal(bytes, &countries)
	if err != nil {
		fmt.Println(err)
	}

	return countries
}

// FindCountry returns the country found in the sentence and if no country is found, returns an empty Country struct
func FindCountry(locale, sentence string) Country {
	for _, country := range countries {
		name, exists := country.Name[locale]

		if !exists {
			continue
		}

		// If the actual country isn't contained in the sentence, continue
		if !strings.Contains(strings.ToLower(sentence), strings.ToLower(name)) {
			continue
		}

		// Returns the right country
		return country
	}

	// Returns an empty country if none has been found
	return Country{}
}
