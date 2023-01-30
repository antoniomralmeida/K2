package inits

import (
	"os"
)

func GetHomeDir() string {
	InitEnvVars()
	if os.Getenv("HOME_DIR") == "" {
		wd, _ := os.Getwd()
		return wd
	} else {
		return os.Getenv("HOME_DIR")
	}
}
