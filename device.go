package natureremo

import "time"

type DeviceCore struct {
	ID                string                     `json:"id"`
	name              string                     `json:"name"`
	TemperatureOffset int64                      `json:"temperature_offset"`
	HumidityOffset    int64                      `json:"humidity_offset"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"created_at"`
	FirmwareVersion   string                     `json:"firmware_version"`
	NewestEvents      map[SensorType]SensorValue `json:"newest_events"`
}

type Device struct {
	DeviceCore
	NewestEvents map[SensorType]SensorValue `json:"newest_events"`
}
