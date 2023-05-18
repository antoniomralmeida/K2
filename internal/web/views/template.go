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
	"github.com/tdewolff/minify/css"
	h "github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

type Template struct {
	FullPath string
	Original string
	Minify   string
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

	T["home_wellcome"] = Minify("text/html", wd+"/web/home_wellcome.gohtml")
	T["login_wellcome"] = Minify("text/html", wd+"/web/login_wellcome.gohtml")
	T["register_wellcome"] = Minify("text/html", wd+"/web/register_wellcome.gohtml")

	T["k2.js"] = Minify("text/javascript", wd+"/web/js/k2.js")
	T["faces.js"] = Minify("text/javascript", wd+"/web/js/faces.js")
	T["olivia.js"] = Minify("text/javascript", wd+"/web/js/olivia.js")
	T["bundle.js"] = Minify("text/javascript", wd+"/web/js/bundle.js")
	T["splash.js"] = Minify("text/javascript", wd+"/web/js/splash.js")
	T["k2.css"] = Minify("text/css", wd+"/web/css/k2.css")

}

func Minify(mediatype string, from string) Template {
	file, err := os.Open(from)
	if err != nil {
		inits.Log(fmt.Sprintf("Error opening file!!! %v", from), inits.Fatal)
		return Template{from, from, from}
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
	m.AddFunc("text/css", css.Minify)

	//m.AddFunc("image/svg+xml", svg.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	err = m.Minify(mediatype, write, read)
	if err != nil {
		inits.Log(err, inits.Error)
		return Template{from, lib.GetFileName(from), lib.GetFileName(from)}
	}
	o.Close()
	f.Close()
	return Template{from, lib.GetFileName(from), lib.GetFileName(to)}
}
