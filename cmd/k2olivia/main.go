package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antoniomralmeida/k2/cmd/k2olivia/locales"
	"github.com/antoniomralmeida/k2/cmd/k2olivia/network"
	"github.com/antoniomralmeida/k2/cmd/k2olivia/server"
	"github.com/antoniomralmeida/k2/cmd/k2olivia/training"
	"github.com/antoniomralmeida/k2/cmd/k2olivia/util"

	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/version"
	"github.com/gookit/color"
)

var neuralNetworks = map[string]network.Network{}

func init() {
	oliviaASCII := string(util.ReadFile("./config/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCII))

	inits.LogInit("k2olivia")
	inits.InitLangs()
	locales.InitStem()
	msg := fmt.Sprintf("Initializing Olivia from K2 KB System, version: %v build: %v PID: %v", version.GetVersion(), version.GetBuild(), os.Getppid())
	fmt.Println(msg)
	fmt.Println("Supported Languages: " + inits.GetSupportedLocales())
	inits.Log(msg, inits.Info)
}

func main() {
	for key := range inits.Locales {
		path := inits.GetHomeDir() + "/k2olivia/res/locales/" + key + "/"
		if ok, _ := lib.Exists(path); !ok {
			err := os.MkdirAll(path, os.ModePerm)
			inits.Log(err, inits.Fatal)
		}
		reTrainModels(key)
	}

	for key := range inits.Locales {
		util.SerializeMessages(key)
		neuralNetworks[key] = training.CreateNeuralNetwork(key, false)
	}

	// Serves the server
	server.Serve(neuralNetworks, os.Getenv("OLIVIA_SERVER_PORT"))
}

// reTrainModels retrain the given locales
func reTrainModels(localesFlag string) {
	// Iterate locales by separating them by comma
	wd := inits.GetHomeDir()
	for _, localeFlag := range strings.Split(localesFlag, ",") {
		path := fmt.Sprintf(wd+"/k2olivia/res/locales/%s/training.json", localeFlag)
		os.Remove(path)
	}
}
