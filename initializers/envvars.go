package initializers

import (
	"github.com/subosito/gotenv"
)

func InitEnvVars() {
	if gotenv.Load("./config/.env") != nil {
		gotenv.Load("../config/.env")
	}
}
