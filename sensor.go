package natureremo

import "time"

// SensorType represents type of sensor.
type SensorType string

const (
	// SensorTypeTemperature represents a temperature sensor.
	SensorTypeTemperature SensorType = "te"
	// SensorTypeHumidity represents a humidity sensor.
	SensorTypeHumidity SensorType = "hu"
	// SensortypeIllumination represents a illumination sensor.
	SensortypeIllumination SensorType = "il"
)

// SensorValue represents value of sensor.
type SensorValue struct {
	Value     float64   `json:"val"`
	CreatedAt time.Time `json:"created_at"`
}
