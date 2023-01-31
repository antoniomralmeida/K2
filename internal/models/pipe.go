package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipe struct {
	id     primitive.ObjectID `json:"_id"`
	avg    float64            `json:"avg"`
	stdDev float64            `json:"stdDev"`
	trust  float64            `json:"trust"`
}
