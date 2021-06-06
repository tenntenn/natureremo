package natureremo

import "time"

// SensorType represents type of sensor.
type SensorType string

const (
	// SensorTypeTemperature represents a temperature sensor.
	SensorTypeTemperature SensorType = "te"
	// SensorTypeHumidity represents a humidity sensor.
	SensorTypeHumidity SensorType = "hu"
	// SensorTypeIllumination represents a illumination sensor.
	SensorTypeIllumination SensorType = "il"
	// SensorTypeMovement represents a movement sensor.
	SensorTypeMovement SensorType = "mo"
)

// SensorValue represents value of sensor.
type SensorValue struct {
	Value     float64   `json:"val"`
	CreatedAt time.Time `json:"created_at"`
}
