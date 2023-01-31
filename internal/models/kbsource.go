package models

type KBSource int8

const (
	Empty KBSource = iota
	User
	IOT
	Simulation
	Inference
)

var KBSourceStr = map[string]KBSource{
	"":           Empty,
	"User":       User,
	"IOT":        IOT,
	"Inference":  Inference,
	"Simulation": Simulation,
}
