package initializers

import (
	"fmt"
	"os"
	"time"

	"github.com/antoniomralmeida/k2/lib"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func LogInit(filebase string) {
	wd, _ := os.Getwd()
	logFileName := wd + os.Getenv("LOGPATH") + filebase + "." + time.Now().Format(lib.YYYYMMDD) + ".json"
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(logFileName + err.Error())
		os.Exit(2)
	}
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.RedirectStdLog(logger)
}

func GetLogger() *zap.Logger {
	return logger
}

const (
	Info  = zapcore.InfoLevel
	Error = zapcore.ErrorLevel
	Fatal = zapcore.FatalLevel
)

func Log(e any, level zapcore.Level) (er error) {
	er = fmt.Errorf("%v", e)
	if e != nil {
		if logger == nil {
			fmt.Println(e)
		} else {
			switch level {
			case zapcore.FatalLevel:
				fmt.Println("Catastrophic error, see log!")
				logger.Fatal(er.Error())
			case zapcore.ErrorLevel:
				logger.Error(er.Error())
			default:
				logger.Info(er.Error())
			}
		}
		return er
	} else {
		return nil
	}
}
