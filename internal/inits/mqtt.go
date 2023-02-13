package inits

import (
	"os"

	"github.com/antoniomralmeida/k2/pkg/dsn"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttClient        mqtt.Client
	messagePubHandler mqtt.MessageHandler
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
	mqttClient = mqtt.NewClient(opts)
}

func SetPubHandler(msgPubHandler mqtt.MessageHandler) {

}

func Publish(topic, json string) {
	if mqttClient != nil {
		token := mqttClient.Publish(topic, 0, false, json)
		token.Wait()
	}
}

func Subscribe(topic string, msgPubHandler mqtt.MessageHandler) {
	messagePubHandler = msgPubHandler
	if mqttClient != nil && messagePubHandler != nil {
		token := mqttClient.Subscribe(topic, 1, nil)
		token.Wait()
	}
}
