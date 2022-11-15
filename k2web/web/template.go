package web

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/tdewolff/minify"
	h "github.com/tdewolff/minify/html"
)

type Template struct {
	original string
	minify   string
}

var T = make(map[string]Template)

func InitTemplates() {
	wd, _ := os.Getwd()
	T["home"] = Minify("text/html", wd+"/k2web/pub/view/template.html")
}

func Minify(mediatype string, from string) Template {
	file, err := os.Open(from)
	if err != nil {
		initializers.Log(fmt.Sprintf("Error opening file!!! %v", from), initializers.Error)
		return Template{from, from}
	}
	defer file.Close()

	o, _ := os.Open(from)
	read := io.Reader(o)
	f, _ := os.CreateTemp("", "tmpfile-")
	write := io.Writer(f)
	to := f.Name()
	m := minify.New()
	//m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", h.Minify)
	//m.AddFunc("text/javascript", js.Minify)
	//m.AddFunc("image/svg+xml", svg.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	err = m.Minify(mediatype, write, read)
	if err != nil {
		initializers.Log(err, initializers.Error)
		return Template{from, from}
	}
	o.Close()
	f.Close()
	nto := to + filepath.Ext(from)
	e := os.Rename(to, nto)
	if e != nil {
		initializers.Log(e, initializers.Error)
	}
	return Template{from, nto}
}
