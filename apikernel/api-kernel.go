package apikernel

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	v "github.com/antoniomralmeida/k2/version"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run(wg *sync.WaitGroup) {
	defer wg.Done()

	app := fiber.New(fiber.Config{AppName: fmt.Sprint("K2 System API-KERNEL ", v.Version, "[", v.Build, "]"),
		DisableStartupMessage: true,
		Prefork:               false})
	wd := initializers.GetHomeDir()
	f := wd + os.Getenv("LOGPATH") + "k2apihttp." + time.Now().Format(lib.YYYYMMDD) + ".log"
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		initializers.Log(fmt.Sprintf("error opening file: %v %v", err, f), initializers.Fatal)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{Output: file,
		TimeFormat: "02/01/2006 15:04:05",
		Format:     "${time} [${ip}:${port}] ${status} ${latency} ${method} ${path} \n"}))

	Routes(app)

	app.Listen(os.Getenv("APIPORT"))
	wg.Done()
}
