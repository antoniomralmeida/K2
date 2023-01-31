package models

type KBSimulation int8

const (
	Default KBSimulation = iota
	MonteCarlo
	NormalDistribution
	LinearRegression
)

var KBSimulationStr = map[string]KBSimulation{
	"":                   Default,
	"MonteCarlo":         MonteCarlo,
	"NormalDistribution": NormalDistribution,
	"LinearRegression":   LinearRegression,
}
