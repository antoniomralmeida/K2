package views

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/gofiber/fiber/v2"

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

	if fiber.IsChild() {
		return
	}
	wd := lib.GetWorkDir()
	T["login"] = Minify("text/html", wd+"/web/login.html")
	T["home"] = Minify("text/html", wd+"/web/home.html")
	T["error"] = Minify("text/html", wd+"/web/error.html")
	T["face"] = Minify("text/html", wd+"/web/face.html")
	T["signup"] = Minify("text/html", wd+"/web/register.html")
	T["splash"] = Minify("text/html", wd+"/web/splash.html")

	T["k2.js"] = Minify("text/javascript", wd+"/web/js/k2.js")
	T["faces.js"] = Minify("text/javascript", wd+"/web/js/faces.js")
	T["olivia.js"] = Minify("text/javascript", wd+"/web/js/olivia.js")
	T["bundle.js"] = Minify("text/javascript", wd+"/web/js/bundle.js")
}

func Minify(mediatype string, from string) Template {
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
	to = strings.TrimSuffix(from, filepath.Ext(from)) + ".min" + filepath.Ext(from)
	f, _ = os.Create(to)
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
		return Template{lib.GetFileName(from), lib.GetFileName(from)}
	}
	o.Close()
	f.Close()
	return Template{lib.GetFileName(from), lib.GetFileName(to)}
}
