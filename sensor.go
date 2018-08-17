package natureremo

import "time"

type SensorType string

const (
	SensorTypeTemperature  SensorType = "te"
	SensorTypeHumidity     SensorType = "hu"
	SensortypeIllumination SensorType = "il"
)

type SensorValue struct {
	Value     float64   `json:"val"`
	CreatedAt time.Time `json:"created_at"`
}
