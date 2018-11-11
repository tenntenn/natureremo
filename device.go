package natureremo

import "time"

// Device represents a device on nature remo.
type Device struct {
	DeviceCore
	NewestEvents map[SensorType]SensorValue `json:"newest_events"`
}

// DeviceCore represents core infomation of a device.
type DeviceCore struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	TemperatureOffset int64     `json:"temperature_offset"`
	HumidityOffset    int64     `json:"humidity_offset"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	FirmwareVersion   string    `json:"firmware_version"`
}
