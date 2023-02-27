package inits

import (
	"log"

	"github.com/subosito/gotenv"
)

func InitEnvVars() {
	if gotenv.Load("./configs/.env") != nil {
		if gotenv.Load("../configs/.env") != nil {
			if err := gotenv.Load("../../configs/.env"); err != nil {
				log.Fatal(err)
			}
		}
	}
}
