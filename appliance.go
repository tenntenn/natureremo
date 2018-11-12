package natureremo

type ApplianceType string

const (
	ApplianceTypeAirConditioner ApplianceType = "AC"
	ApplianceTypeIR             ApplianceType = "IR"
)

type ApplianceModel struct {
	ID           string `json:"id"`
	Manufacturer string `json:"manufacturer"`
	RemoteName   string `json:"remote_name"`
	Name         string `json:"name"`
	Image        string `json:"image"`
}

type Appliance struct {
	ID             string          `json:"id"`
	Device         *DeviceCore     `json:"device"`
	Model          *ApplianceModel `json:"model"`
	Nickname       string          `json:"nickname"`
	Image          string          `json:"image"`
	Type           ApplianceType   `json:"type"`
	AirConSettings *AirConParams   `json:"settings"`
	AirCon         *AirCon         `json:"aircon"`
	Signals        []*Signal       `json:"signals"`
}
