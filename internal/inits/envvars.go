package inits

import (
	"github.com/subosito/gotenv"
)

func InitEnvVars() {
	if gotenv.Load("./configs/.env") != nil {
		gotenv.Load("../configs/.env")
	}
}
