package initializers

import (
	"log"
	"os"

	"github.com/leemcloughlin/logfile"
)

func LogInit() {

	logFileName := os.Getenv("LOG")
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
