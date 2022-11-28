package initializers

import (
	"log"

	"github.com/eduardonunesp/goaiml"
)

var aiml map[string]*goaiml.AIML

func InitAiml(uid string) {
	if aiml == nil {
		aiml = make(map[string]*goaiml.AIML)
	}
	if aiml[uid] != nil {
		return
	}
	aiml[uid] = goaiml.NewAIML()
	err := aiml[uid].Learn("./config/k2.aiml.xml")
	if err != nil {
		log.Fatal(err)
	}
}

func GetResponse(uid, line string) string {
	InitAiml(uid)
	ret, _ := aiml[uid].Respond(line)
	return ret
}
