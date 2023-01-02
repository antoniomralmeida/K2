package iot

import (
	"fmt"
	"os"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	v "github.com/antoniomralmeida/k2/version"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Run() {
	app := fiber.New(fiber.Config{AppName: fmt.Sprint("K2 IoT MidWare ", v.Version, "[", v.Build, "]"),
		DisableStartupMessage: true,
		Prefork:               false})

	wd := initializers.GetHomeDir()
	f := wd + os.Getenv("LOGPATH") + "k2iothttp." + time.Now().Format(lib.YYYYMMDD) + ".log"
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		initializers.Log(fmt.Sprintf("error opening file: %v %v", err, f), initializers.Fatal)
	}
	defer file.Close()

	app.Use(logger.New(logger.Config{Output: file,
		TimeFormat: "02/01/2006 15:04:05",
		Format:     "${time} [${ip}:${port}] ${status} ${latency} ${method} ${path} \n"}))
	app.Use(requestid.New())
	// Provide a custom compression level
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	Routes(app)

	app.Listen(os.Getenv("APIPORT"))
}
