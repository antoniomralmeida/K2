package inits

import (
	"fmt"
	"os"

	"github.com/antoniomralmeida/k2/pkg/dsn"
	"github.com/olivia-ai/olivia-kit-go/chat"
)

var olivia chat.Client

type Configuration struct {
	Port       string `json:"port"`
	Host       string `json:"host"`
	SSL        bool   `json:"ssl"`
	DebugLevel string `json:"debug_level"`
	BotName    string `json:"bot_name"`
	UserToken  string `json:"user_token"`
}

func InitOlivia() {
	//Init client Olivia
	server := os.Getenv("OLIVIA_SERVER")
	dsn := dsn.Decode(server)
	config := Configuration{Host: dsn.Host(), Port: dsn.Port(), SSL: false, DebugLevel: "error", BotName: dsn.Query("botname")}
	fmt.Println(config)
	var information map[string]interface{}
	client, err := chat.NewClient(
		fmt.Sprintf("%s:%s", config.Host, config.Port),
		config.SSL,
		&information,
	)
	Log(err, Fatal)
	defer client.Close()
	olivia = client
}

func GetResponse(locale, uid, msg string) string {
	olivia.Locale = locale
	olivia.Token = uid
	response, err := olivia.SendMessage(msg)
	Log(err, Error)
	return response.Content
}
