package iot

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Protocol string

const (
	XMPP   Protocol = "XMPP"
	MQTT   Protocol = "MQTT"
	CoAP   Protocol = "CoAP"
	SOAP   Protocol = "SOAP"
	HTTP   Protocol = "HTTP"
	SNMP   Protocol = "SNMP"
	MODBUS Protocol = "MODBUS"
)

var ProtocolName = map[Protocol]string{XMPP: "Extensible Messaging and Presence Protocol",
	MQTT:   "Message Queuing Telemetry Transport",
	CoAP:   "Constrained Application Protocol",
	SOAP:   "Simple Object Access Protocol",
	HTTP:   "Hypertext Transfer Protocol",
	SNMP:   "Simple Network Management Protocol",
	MODBUS: "Data communications protocol for PLC",
}

type Access byte

const (
	ReadOnly Access = 1
	ReadWrite
)

type Gate struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Protocol Protocol           `bson:"protocol" json:"protocol"`
	Access   Access             `bson:"access" json:"access"`
	Url      string             `bson:"url" json:"url"`
	Options  string             `bson:"options" json"options"`
}

func (me Protocol) String() string {
	return ProtocolName[me]
}
