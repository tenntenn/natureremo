package natureremo

import "time"

// Device represents a device such as Nature Remo and Nature Remo Mini.
type Device struct {
	DeviceCore
	// NewestEvents is newest sensor values such as temperature, humidity and illumination.
	NewestEvents map[SensorType]SensorValue `json:"newest_events"`
}

// DeviceCore represents core information of a device.
type DeviceCore struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	TemperatureOffset int64     `json:"temperature_offset"`
	HumidityOffset    int64     `json:"humidity_offset"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	FirmwareVersion   string    `json:"firmware_version"`
	MacAddress        string    `json:"mac_address"`
	BtMacAddress      string    `json:"bt_mac_address"`
	SerialNumber      string    `json:"serial_number"`
	Users             []User    `json:"users"`
}
