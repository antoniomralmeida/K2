package web

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/antoniomralmeida/k2/inits"
	"github.com/tdewolff/minify"
	h "github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

type Template struct {
	original string
	minify   string
}

var T = make(map[string]Template)

func InitTemplates() {

	wd := inits.GetHomeDir()
	T["login"] = Minify("text/html", wd+"/k2web/pub/login.html", true)
	T["home"] = Minify("text/html", wd+"/k2web/pub/home.html", true)
	T["404"] = Minify("text/html", wd+"/k2web/pub/404.html", true)
	T["face"] = Minify("text/html", wd+"/k2web/pub/face.html", true)
	T["sigup"] = Minify("text/html", wd+"/k2web/pub/register.html", true)
	T["k2.js"] = Minify("text/javascript", wd+"/k2web/pub/js/k2.js", false)
	T["faces.js"] = Minify("text/javascript", wd+"/k2web/pub/js/faces.js", false)
	T["olivia.js"] = Minify("text/javascript", wd+"/k2web/pub/js/olivia.js", false)
	T["bundle.js"] = Minify("text/javascript", wd+"/k2web/pub/js/bundle.js", false)
}

func Minify(mediatype string, from string, temp bool) Template {
	file, err := os.Open(from)
	if err != nil {
		inits.Log(fmt.Sprintf("Error opening file!!! %v", from), inits.Fatal)
		return Template{from, from}
	}
	defer file.Close()

	o, _ := os.Open(from)
	read := io.Reader(o)
	var to string
	var f *os.File
	if temp {
		f, _ = os.CreateTemp("", "tmpfile-")
		to = f.Name()
	} else {
		to = strings.TrimSuffix(from, filepath.Ext(from)) + ".min" + filepath.Ext(from)
		f, _ = os.Create(to)
	}
	write := io.Writer(f)

	m := minify.New()
	//m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", h.Minify)
	m.AddFunc("text/javascript", js.Minify)
	//m.AddFunc("image/svg+xml", svg.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	err = m.Minify(mediatype, write, read)
	if err != nil {
		inits.Log(err, inits.Error)
		return Template{from, from}
	}
	o.Close()
	f.Close()
	if temp {
		nto := to + filepath.Ext(from)
		e := os.Rename(to, nto)
		if e != nil {
			inits.Log(e, inits.Error)
		}
		return Template{from, nto}
	} else {
		return Template{from, to}
	}
}
