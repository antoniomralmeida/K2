package lib

import (
	"bufio"
	"os"
	"strings"
)

type Stem struct {
	locale string
	stams  map[string]string
}

func (s *Stem) Stem(stem string) string {
	return s.stams[stem]
}

func NewStem(locale string) (*Stem, error) {
	stam := new(Stem)
	stam.locale = locale
	stam.stams = make(map[string]string)

	err, filename := DownloadFile("https://raw.githubusercontent.com/michmech/lemmatization-lists/master/lemmatization-"+locale+".txt", "./k2olivia/res/locales/"+locale+"/")
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
			stam.stams[fields[0]] = fields[1]
		}
	}

	readFile.Close()
	return stam, nil
}
