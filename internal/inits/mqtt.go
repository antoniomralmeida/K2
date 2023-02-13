package inits

import (
	"os"

	"github.com/antoniomralmeida/k2/pkg/dsn"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func InitMQTT() {
	uri := os.Getenv("MQTT_SERVER")
	parts := dsn.Decode(uri)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(parts.Short())
	opts.SetClientID("k2_mqtt_client")
	opts.SetUsername(parts.User())
	opts.SetPassword(parts.Password())
	opts.SetDefaultPublishHandler(messagePubHandler)

	client := mqtt.NewClient(opts)

}
