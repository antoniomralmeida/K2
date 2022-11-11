package initializers

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

const YYYYMMDD = "2006-01-02"

func LogInit(filebase string) {
	wd, _ := os.Getwd()
	logFileName := wd + os.Getenv("LOGPATH") + filebase + "." + time.Now().Format(YYYYMMDD) + ".json"
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

const (
	Info  = zapcore.InfoLevel
	Error = zapcore.ErrorLevel
	Fatal = zapcore.FatalLevel
)

func Log(e any, level zapcore.Level) (er error) {
	er = fmt.Errorf("%v", e)
	if e != nil {
		switch level {
		case zapcore.FatalLevel:
			logger.Fatal(er.Error())
			fmt.Println("Catastrophic error, see log!")
		case zapcore.ErrorLevel:
			logger.Error(er.Error())
		default:
			logger.Info(er.Error())
		}
		return er
	} else {
		return nil
	}
}
