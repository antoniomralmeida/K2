package models

type KBSource int8

const (
	Undefined KBSource = iota
	FromUser
	IOT
	Simulation
	Inference
)

var KBSourceStr = map[string]KBSource{
	"":           Undefined,
	"User":       FromUser,
	"IOT":        IOT,
	"Inference":  Inference,
	"Simulation": Simulation,
}

func ToKBSources(sources []string) (ret []KBSource) {
	for i := range sources {
		ret = append(ret, KBSourceStr[sources[i]])
	}
	return
}
