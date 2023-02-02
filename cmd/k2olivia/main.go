package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/internal/olivia/locales"
	"github.com/antoniomralmeida/k2/internal/olivia/network"
	"github.com/antoniomralmeida/k2/internal/olivia/server"
	"github.com/antoniomralmeida/k2/internal/olivia/training"
	"github.com/antoniomralmeida/k2/internal/olivia/util"

	"github.com/gookit/color"
)

var neuralNetworks = map[string]network.Network{}

func init() {
	oliviaASCII := string(util.ReadFile("./configs/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCII))

	inits.LogInit("k2olivia")
	inits.InitLangs()
	locales.InitStem()
	msg := fmt.Sprintf("Initializing Olivia from K2 KB System, version: %v build: %v PID: %v", lib.GetVersion(), lib.GetBuild(), os.Getppid())
	fmt.Println(msg)
	fmt.Println("Supported Languages: " + inits.GetSupportedLocales())
	inits.Log(msg, inits.Info)
}

func main() {
	for key := range inits.Locales {
		path := inits.GetHomeDir() + "/data/locales/" + key + "/"
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
		path := fmt.Sprintf(wd+"/data/locales/%s/training.json", localeFlag)
		os.Remove(path)
	}
}
