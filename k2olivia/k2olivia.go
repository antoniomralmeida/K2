package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/k2olivia/locales"
	"github.com/antoniomralmeida/k2/k2olivia/network"
	"github.com/antoniomralmeida/k2/k2olivia/server"
	"github.com/antoniomralmeida/k2/k2olivia/training"
	"github.com/antoniomralmeida/k2/k2olivia/util"
	"github.com/antoniomralmeida/k2/version"
	"github.com/gookit/color"
)

var neuralNetworks = map[string]network.Network{}

func init() {
	initializers.LogInit("k2olivia")
<<<<<<< HEAD
	msg := fmt.Sprintf("Initializing K2 Olivia version: %v build: %v PID: %v", version.GetVersion(), version.GetBuild(), os.Getppid())
=======
	initializers.InitLangs()
	locales.InitStem()
	msg := fmt.Sprintf("Initializing K2 Olivia version: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
>>>>>>> 01887a253f097f28bcbfe9116bed04d1b593fab3
	fmt.Println(msg)
	initializers.Log(msg, initializers.Info)
}

func main() {
	for key := range initializers.Locales {
		reTrainModels(key)
	}
	wd := initializers.GetHomeDir()
	// Print the Olivia ascii text
	oliviaASCII := string(util.ReadFile(wd + "/k2olivia/res/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCII))

	for key := range initializers.Locales {
		util.SerializeMessages(key)
		neuralNetworks[key] = training.CreateNeuralNetwork(key, false)
	}

	// Serves the server
	server.Serve(neuralNetworks, os.Getenv("OLIVIA_SERVER_PORT"))
}

// reTrainModels retrain the given locales
func reTrainModels(localesFlag string) {
	// Iterate locales by separating them by comma
	wd := initializers.GetHomeDir()
	for _, localeFlag := range strings.Split(localesFlag, ",") {
		path := fmt.Sprintf(wd+"/k2olivia/res/locales/%s/training.json", localeFlag)
		os.Remove(path)
	}
}
