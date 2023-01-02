package natureremo

// Appliance represents controlable devices with Nature Remo
// such as air conditioners and TVs.
type Appliance struct {
	ID             string          `json:"id"`
	Type           ApplianceType   `json:"type"`
	Device         *DeviceCore     `json:"device"`
	Model          *ApplianceModel `json:"model"`
	Nickname       string          `json:"nickname"`
	Image          string          `json:"image"`
	Signals        []*Signal       `json:"signals"`
	AirConSettings *AirConSettings `json:"settings"`
	AirCon         *AirCon         `json:"aircon"`
	TV             *TV             `json:"tv"`
	Light          *Light          `json:"light"`
	SmartMeter     *SmartMeter     `json:"smart_meter"`
}

// SignalByName gets a signal by name from Signals.
// If there are not any signals which have specified name,
// SignalByName returns nil.
func (a *Appliance) SignalByName(name string) *Signal {
	for _, s := range a.Signals {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// ApplianceType represents type of appliance.
type ApplianceType string

const (
	// ApplianceTypeAirCon represents an air conditioner.
	ApplianceTypeAirCon ApplianceType = "AC"
	// ApplianceTypeTV represents an TV
	ApplianceTypeTV ApplianceType = "TV"
	// ApplianceTypeLight represents Light
	ApplianceTypeLight ApplianceType = "LIGHT"
	// ApplianceTypeIR represents a device which is controlled by infrared.
	ApplianceTypeIR ApplianceType = "IR"
	// ApplianceTypeSmartMeter represents a smart meter
	ApplianceTypeSmartMeter ApplianceType = "EL_SMART_METER"
)

// ApplianceModel is device information of appliance
// which is identified by Nature Remo API.
type ApplianceModel struct {
	ID string `json:"id"`
	// Not in swagger, but actually included in the response
	Country      string `json:"country"`
	Manufacturer string `json:"manufacturer"`
	RemoteName   string `json:"remote_name"`
	Name         string `json:"name"`
	Image        string `json:"image"`
}
