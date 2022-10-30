package lib

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leemcloughlin/logfile"
)

func IsNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 32)
	return err == nil
}

func LogInit() {

	logFileName := "./log/k2.log"
	logFile, err := logfile.New(
		&logfile.LogFile{
			FileName: logFileName,
			MaxSize:  500 * 1024, // 500K duh!
			Flags:    logfile.FileOnly | logfile.OverWriteOnStart})
	if err != nil {
		log.Fatalf("Failed to create logFile %s: %s\n", logFileName, err)
	}

	log.SetOutput(logFile)
}

func LogFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Log(msg string) {
	log.Println(msg)
}

func IsMainThread() bool {
	return !fiber.IsChild()
}
