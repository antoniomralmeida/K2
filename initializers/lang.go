package initializers

import (
	"bufio"
	"os"
	"strings"

	"github.com/antoniomralmeida/k2/lib"
)

type Locale struct {
	Description       string
	SpeechSynthesisId int
	Stemmer           *Stem
}

const DefaultLocale = "en"

var Locales map[string]Locale

func inLocales(locale string) bool {
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
		ret = ret + " " + key + ":" + value.Description
	}
	return
}

func InitLangs() {
	Locales = make(map[string]Locale)
	Locales["en"] = Locale{Description: "English", SpeechSynthesisId: 1}
	Locales["pt"] = Locale{Description: "Portuguese(BR)", SpeechSynthesisId: 0}
	if inLocales("es") {
		Locales["es"] = Locale{Description: "Spanish", SpeechSynthesisId: 262}
	}
	if inLocales("de") {
		Locales["de"] = Locale{Description: "German", SpeechSynthesisId: 143}
	}
	if inLocales("hi") {
		Locales["hi"] = Locale{Description: "Hindi", SpeechSynthesisId: 154}
	}
	if inLocales("ar") {
		Locales["ar"] = Locale{Description: "Arabic", SpeechSynthesisId: 12}
	}
	if inLocales("de") {
		Locales["bn"] = Locale{Description: "Bengali", SpeechSynthesisId: 48}
	}
	if inLocales("ru") {
		Locales["ru"] = Locale{Description: "Russian", SpeechSynthesisId: 213}
	}
	if inLocales("ja") {
		Locales["ja"] = Locale{Description: "Japanese", SpeechSynthesisId: 167}
	}
	if inLocales("fr") {
		Locales["fr"] = Locale{Description: "French", SpeechSynthesisId: 133}
	}
	if inLocales("it") {
		Locales["it"] = Locale{Description: "Italian", SpeechSynthesisId: 164}
	}
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
