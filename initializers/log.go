package initializers

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func LogInit() {
	wd, _ := os.Getwd()

	logFileName := wd + os.Getenv("K2LOG")
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
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
