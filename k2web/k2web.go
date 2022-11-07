package main

import (
	"log"
	"os"

	"github.com/antoniomralmeida/k2/k2web/web"
	"github.com/leemcloughlin/logfile"
	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load(web.GetK2Path() + "/.env"); err != nil {
		log.Fatal(err)
	}
	logFileName := web.GetK2Path() + os.Getenv("LOGWEB")
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
func main() {
	web.Run()
}
