package inits

import (
	"os"
)

func GetHomeDir() string {
	InitEnvVars()
	if os.Getenv("WORKDIR") == "" {
		wd, _ := os.Getwd()
		return wd
	} else {
		return os.Getenv("WORKDIR")
	}
}
