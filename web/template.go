package web

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify"
	h "github.com/tdewolff/minify/html"
)

type Template struct {
	original string
	minify   string
}

var T = make(map[string]Template)

func InitTemplates() {
	T["home"] = Minify("text/html", "./web/pub/view/template.html")
}

func Minify(mediatype string, from string) Template {
	file, err := os.Open(from)
	if err != nil {
		fmt.Println("Error opening file!!!")
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
		panic(err)
	}

	o.Close()
	f.Close()
	nto := to + filepath.Ext(from)
	e := os.Rename(to, nto)
	if e != nil {
		log.Fatal(e)
	}
	return Template{from, nto}
}
