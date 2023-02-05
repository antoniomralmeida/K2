package golibretranslate

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type sentence struct {
	Translate string `json:"trans"`
	Orignal   string `json:"orig"`
}
type response struct {
	Sentences []sentence
	Source    string `json:"src"`
}

func (r *response) getTranslate() string {
	if len(r.Sentences) > 0 {
		return r.Sentences[0].Translate
	}
	return ""
}

func Translate(text string, locale_source, locale_target string) (string, error) {

	data := url.Values{
		"q":  {text},
		"sl": {locale_source},
		"tl": {locale_target},
	}

	u, _ := url.Parse("https://translate.google.com/translate_a/single")
	q := u.Query()
	q.Set("client", "at")
	q.Set("dt", "t")
	q.Set("dj", "1")
	u.RawQuery = q.Encode()

	resp, err := http.PostForm(u.String(), data)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	res := new(response)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err == nil {
		return res.getTranslate(), nil
	} else {
		return "", err
	}
}
