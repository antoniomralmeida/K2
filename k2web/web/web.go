package web

import (
	"fmt"
	"os"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/models"
	"github.com/antoniomralmeida/k2/version"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var ctxweb = models.Context{}

func Run() {
	InitLangs()
	InitTemplates()
	initializers.ConnectDB()

	ctxweb.ApiKernel = os.Getenv("APIKERNEL")
	ctxweb.Avatar = os.Getenv("AVATAR")
	app := fiber.New(fiber.Config{AppName: fmt.Sprint("K2 System ", version.Version, "[", version.Build, "]"),
		DisableStartupMessage: false,
		Prefork:               true})
	wd, _ := os.Getwd()
	f := wd + os.Getenv("LOGPATH") + "k2webhttp." + time.Now().Format(lib.YYYYMMDD) + ".log"
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err, f)
		initializers.Log(fmt.Sprintf("error opening file: %v %v", err, f), initializers.Fatal)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{Output: file,
		TimeFormat: "02/01/2006 15:04:05",
		Format:     "${time} [${ip}:${port}] ${status} ${latency} ${method} ${path} \n"}))

	Routes(app)
	app.Listen(os.Getenv("HTTPPORT"))
}
