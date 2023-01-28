package initializers

import (
	"bufio"
	"os"
	"strings"

	"github.com/antoniomralmeida/k2/lib"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type Locale struct {
	Description       string
	SpeechSynthesisId int
	Stemmer           *Stem
}

var DefaultLocale = language.English.String()

var Locales map[string]Locale

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
func NewSupportedLanguage(locale string, SynthesisId int) {
	if inLocalesConfig(locale) || locale == language.English.String() || locale == language.Portuguese.String() {
		toen := display.English.Languages()
		tag := language.MustParse(locale)
		Locales[locale] = Locale{Description: toen.Name(tag), SpeechSynthesisId: SynthesisId}
	}
}

func InitLangs() {
	Locales = make(map[string]Locale)
	NewSupportedLanguage(language.English.String(), 1)
	NewSupportedLanguage(language.Portuguese.String(), 0)
	NewSupportedLanguage(language.Spanish.String(), 262)
	NewSupportedLanguage(language.German.String(), 143)
	NewSupportedLanguage(language.Hindi.String(), 154)
	NewSupportedLanguage(language.Arabic.String(), 12)
	NewSupportedLanguage(language.Bengali.String(), 48)
	NewSupportedLanguage(language.Russian.String(), 213)
	NewSupportedLanguage(language.Japanese.String(), 167)
	NewSupportedLanguage(language.French.String(), 133)
	NewSupportedLanguage(language.Italian.String(), 164)
	NewSupportedLanguage(language.Chinese.String(), 68)

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

	err, filename := lib.DownloadFile("https://raw.githubusercontent.com/michmech/lemmatization-lists/master/lemmatization-"+locale+".txt", GetHomeDir()+"/k2olivia/res/locales/"+locale+"/")
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
