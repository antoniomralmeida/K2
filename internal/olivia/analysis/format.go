package analysis

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/antoniomralmeida/golibretranslate"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/olivia/locales"
	"github.com/antoniomralmeida/k2/internal/olivia/util"
)

// arrange checks the format of a string to normalize it, remove ignored characters
func (sentence *Sentence) arrange() {
	// Remove punctuation after letters
	punctuationRegex := regexp.MustCompile(`[a-zA-Z]( )?(\.|\?|!|¿|¡)`)
	sentence.Content = punctuationRegex.ReplaceAllStringFunc(sentence.Content, func(s string) string {
		punctuation := regexp.MustCompile(`(\.|\?|!)`)
		return punctuation.ReplaceAllString(s, "")
	})

	sentence.Content = strings.ReplaceAll(sentence.Content, "-", " ")
	sentence.Content = strings.TrimSpace(sentence.Content)
}
func stopWorksFile(locale string) string {
	return inits.GetHomeDir() + "/data/locales/" + locale + "/stopwords.txt"
}

// removeStopWords takes an arary of words, removes the stopwords and returns it
func removeStopWords(locale string, words []string) []string {
	// Don't remove stopwords for small sentences like “How are you” because it will remove all the words
	if len(words) <= 4 {
		return words
	}

	stopfile := stopWorksFile(locale)
	if ok, _ := lib.Exists(stopfile); !ok {
		tmpFile := intentsFile(inits.DefaultLocale)
		bytes, _ := ioutil.ReadFile(tmpFile)
		tmpWords := string(bytes)
		tmpWords, err := golibretranslate.Translate(tmpWords, inits.DefaultLocale, locale)
		inits.Log(err, inits.Error)
		f, err := os.Create(stopfile)
		inits.Log(err, inits.Error)
		f.WriteString(tmpWords)
		f.Close()
	}

	// Read the content of the stopwords file
	txt, _ := ioutil.ReadFile(stopfile)
	stopWords := string(txt)

	var wordsToRemove []string

	// Iterate through all the stopwords
	for _, stopWord := range strings.Split(stopWords, "\n") {
		// Iterate through all the words of the given array
		for _, word := range words {
			// Continue if the word isn't a stopword
			if !strings.Contains(stopWord, word) {
				continue
			}

			wordsToRemove = append(wordsToRemove, word)
		}
	}

	return util.Difference(words, wordsToRemove)
}

// tokenize returns a list of words that have been lower-cased
func (sentence Sentence) tokenize() (tokens []string) {
	// Split the sentence in words
	tokens = strings.Fields(sentence.Content)

	// Lower case each word
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}

	tokens = removeStopWords(sentence.Locale, tokens)

	return
}

// stem returns the sentence split in stemmed words
func (sentence Sentence) stem() (tokenizeWords []string) {
	locale := locales.GetLocaleByName(sentence.Locale)

	tokens := sentence.tokenize()

	// Get the string token and push it to tokenizeWord
	if locale.Stemmer != nil {
		for _, tokenizeWord := range tokens {
			word := locale.Stemmer.Stem(tokenizeWord)
			tokenizeWords = append(tokenizeWords, word...)
		}
	}
	return
}

// WordsBag retrieves the intents words and returns the sentence converted in a bag of words
func (sentence Sentence) WordsBag(words []string) (bag []float64) {
	for _, word := range words {
		// Append 1 if the patternWords contains the actual word, else 0
		var valueToAppend float64
		if util.Contains(sentence.stem(), word) {
			valueToAppend = 1
		}

		bag = append(bag, valueToAppend)
	}

	return bag
}
