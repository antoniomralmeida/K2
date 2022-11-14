package initializers

import (
	"github.com/subosito/gotenv"
)

func InitEnvVars() {
	gotenv.Load("./bin/.env")
}
