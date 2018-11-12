package natureremo

type DetectedAircon struct {
	Model  *ApplianceModel `json:"model"`
	Params *AirConParams   `json:"params"`
}

type AirConParams struct {
	Temperature   string        `json:"temp"`
	OperationMode OperationMode `json:"mode"`
	AirVolume     AirVolume     `json:"vol"`
	AirDirection  AirDirection  `json:"dir"`
	Button        Button        `json:"button"`
}

type OperationMode string

func (om OperationMode) StringValue() string {
	return string(om)
}

const (
	OperationModeAuto OperationMode = "auto"
	OperationModeCool OperationMode = "cool"
	OperationModeWarm OperationMode = "warm"
	OperationModeDry  OperationMode = "dry"
	OperationModeBlow OperationMode = "blow"
)

type AirConRangeMode struct {
	Temperature  []string       `json:"temp"`
	AirVolume    []AirVolume    `json:"vol"`
	AirDirection []AirDirection `json:"dir"`
}

type AirVolume string

func (v AirVolume) StringValue() string {
	return string(v)
}

const (
	AirVolumeAuto = ""
	AirVolume1    = "1"
	AirVolume2    = "2"
	AirVolume3    = "3"
	AirVolume4    = "4"
	AirVolume5    = "5"
	AirVolume6    = "6"
	AirVolume7    = "7"
	AirVolume8    = "8"
	AirVolume9    = "9"
	AirVolume10   = "10"
)

type AirDirection string

func (d AirDirection) StringValue() string {
	return string(d)
}

const (
	AirDirectionAuto AirDirection = ""
)

type Button string

func (b Button) StringValue() string {
	return string(b)
}

const (
	ButtonPowerOn  Button = ""
	ButtonPowerOff Button = "power-off"
)

type AirCon struct {
	Range           *Range          `json:"range"`
	TemperatureUnit TemperatureUnit `json:"tempUnit"`
}

type Range struct {
	AirConRangeModes *AirConRangeModes `json:"modes"`
	FixedButtons     []Button          `json:"fixedButtons"`
}

type AirConRangeModes struct {
	Cool *AirConRangeMode `json:"cool"`
	Warm *AirConRangeMode `json:"warm"`
	Dry  *AirConRangeMode `json:"dry"`
	Blow *AirConRangeMode `json:"blow"`
	Auto *AirConRangeMode `json:"auto"`
}

type TemperatureUnit string

const (
	TemperatureUnitAuto       TemperatureUnit = ""
	TemperatureUnitFahrenheit TemperatureUnit = "f"
	TemperatureUnitCelsius    TemperatureUnit = "c"
)
