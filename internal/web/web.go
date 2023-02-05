package web

import (
	"fmt"
	"os"
	"time"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/web/context"
	"github.com/antoniomralmeida/k2/internal/web/controllers"
	"github.com/antoniomralmeida/k2/internal/web/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run() {
	context.InitLangs()
	views.InitTemplates()
	inits.ConnectDB()

	context.Ctxweb.ApiKernel = os.Getenv("APIKERNEL")
	context.Ctxweb.Avatar = os.Getenv("AVATAR")
	app := fiber.New(fiber.Config{AppName: fmt.Sprint("K2 KB System ", lib.GetVersion(), "[", lib.GetBuild(), "]"),
		DisableStartupMessage: false,
		Prefork:               false})
	wd := inits.GetHomeDir()
	f := wd + os.Getenv("LOGPATH") + "k2webhttp." + time.Now().Format(lib.YYYYMMDD) + ".log"
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err, f)
		inits.Log(fmt.Sprintf("error opening file: %v %v", err, f), inits.Fatal)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{Output: file,
		TimeFormat: "02/01/2006 15:04:05",
		Format:     "${time} [${ip}:${port}] ${status} ${latency} ${method} ${path} \n"}))

	controllers.Routes(app)
	app.Listen(os.Getenv("HTTPPORT"))
}
