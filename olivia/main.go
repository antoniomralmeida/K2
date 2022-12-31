package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/olivia/dashboard"
	"github.com/antoniomralmeida/k2/olivia/locales"
	"github.com/antoniomralmeida/k2/olivia/network"
	"github.com/antoniomralmeida/k2/olivia/server"
	"github.com/antoniomralmeida/k2/olivia/training"
	"github.com/antoniomralmeida/k2/olivia/util"
	"github.com/antoniomralmeida/k2/version"

	"github.com/gookit/color"
)

var neuralNetworks = map[string]network.Network{}

func init() {
	initializers.InitEnvVars()
	initializers.LogInit("k2olivia")
	msg := fmt.Sprintf("Initializing K2 Olivia version: %v build: %v PID: %v", version.Version, version.Build, os.Getppid())
	fmt.Println(msg)
	initializers.Log(msg, initializers.Info)
}

func main() {
	port := flag.String("port", "8090", "The port for the API and WebSocket.")
	localesFlag := flag.String("re-train", "", "The locale(s) to re-train.")
	flag.Parse()

	// If the locales flag isn't empty then retrain the given models
	if *localesFlag != "" {
		reTrainModels(*localesFlag)
	}

	// Print the Olivia ascii text
	oliviaASCII := string(util.ReadFile("./olivia/res/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCII))

	// Create the authentication token
	dashboard.Authenticate()

	for _, locale := range locales.Locales {
		util.SerializeMessages(locale.Tag)

		neuralNetworks[locale.Tag] = training.CreateNeuralNetwork(
			locale.Tag,
			false,
		)
	}

	// Get port from environment variables if there is
	if os.Getenv("OLIVIA_SERVER_PORT") != "" {
		*port = os.Getenv("OLIVIA_SERVER_PORT")
	}

	// Serves the server
	server.Serve(neuralNetworks, *port)
}

// reTrainModels retrain the given locales
func reTrainModels(localesFlag string) {
	// Iterate locales by separating them by comma
	for _, localeFlag := range strings.Split(localesFlag, ",") {
		path := fmt.Sprintf("./olivia/res/locales/%s/training.json", localeFlag)
		err := os.Remove(path)

		if err != nil {
			fmt.Printf("Cannot re-train %s model.", localeFlag)
			return
		}
	}
}
