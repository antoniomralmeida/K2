package inits

import (
	"log"

	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/subosito/gotenv"
)

func InitEnvVars() {
	if err := gotenv.Load(lib.GetWorkDir() + "/configs/.env"); err != nil {
		log.Fatal(err)
	}
}
