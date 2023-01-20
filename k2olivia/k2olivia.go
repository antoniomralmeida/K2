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
	msg := fmt.Sprintf("Initializing K2 Olivia version: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
	fmt.Println(msg)
	initializers.Log(msg, initializers.Info)
}

func main() {
	for _, locale := range locales.Locales {
		reTrainModels(locale.Tag)
	}
	wd := initializers.GetHomeDir()
	// Print the Olivia ascii text
	oliviaASCII := string(util.ReadFile(wd + "/k2olivia/res/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCII))

	// Create the authentication token

	//dashboard.Authenticate()
	for _, locale := range locales.Locales {
		util.SerializeMessages(locale.Tag)
		neuralNetworks[locale.Tag] = training.CreateNeuralNetwork(locale.Tag, false)
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
