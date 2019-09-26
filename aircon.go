package natureremo

// AirCon represents air conditioner.
type AirCon struct {
	Range           *AirConRange    `json:"range"`
	TemperatureUnit TemperatureUnit `json:"tempUnit"`
}

// AirConRange holds range of each setting and buttons of air conditioner.
type AirConRange struct {
	Modes        map[OperationMode]*AirConRangeMode `json:"modes"`
	FixedButtons []Button                           `json:"fixedButtons"`
}

// TemperatureUnit is unit of temperature.
type TemperatureUnit string

const (
	// TemperatureUnitAuto represents auto.
	TemperatureUnitAuto TemperatureUnit = ""
	// TemperatureUnitFahrenheit represents fahrenheit.
	TemperatureUnitFahrenheit TemperatureUnit = "f"
	// TemperatureUnitCelsius represents celsius.
	TemperatureUnitCelsius TemperatureUnit = "c"
)

// DetectedAirCon represents information of an air conditioner
// which is detected by ApplianceService.Detected.
type DetectedAirCon struct {
	Model  *ApplianceModel `json:"model"`
	Params *AirConSettings `json:"params"`
}

// AirConSettings represents settings of air conditioner.
type AirConSettings struct {
	Temperature   string        `json:"temp"`
	OperationMode OperationMode `json:"mode"`
	AirVolume     AirVolume     `json:"vol"`
	AirDirection  AirDirection  `json:"dir"`
	Button        Button        `json:"button"`
}

// OperationMode represents a operation mode of air conditioner such as warm and cool.
type OperationMode string

// StringValue converts OperationMode to string value.
func (om OperationMode) StringValue() string {
	return string(om)
}

const (
	// OperationModeAuto represents auto mode.
	OperationModeAuto OperationMode = "auto"
	// OperationModeCool represents cool mode.
	OperationModeCool OperationMode = "cool"
	// OperationModeWarm represents warm mode.
	OperationModeWarm OperationMode = "warm"
	// OperationModeDry represents dry mode.
	OperationModeDry OperationMode = "dry"
	// OperationModeBlow represents blow mode.
	OperationModeBlow OperationMode = "blow"
)

// AirConRangeMode represetns ranges of settings of air conditioner.
type AirConRangeMode struct {
	Temperature  []string       `json:"temp"`
	AirVolume    []AirVolume    `json:"vol"`
	AirDirection []AirDirection `json:"dir"`
}

// AirVolume represents air volume of air conditioner.
type AirVolume string

// StringValue converts AirVolume to string.
func (v AirVolume) StringValue() string {
	return string(v)
}

const (
	// AirVolumeAuto represents auto.
	AirVolumeAuto = "auto"
	// AirVolume1 represents volume 1.
	AirVolume1 = "1"
	// AirVolume2 represents volume 2.
	AirVolume2 = "2"
	// AirVolume3 represents volume 3.
	AirVolume3 = "3"
	// AirVolume4 represents volume 4.
	AirVolume4 = "4"
	// AirVolume5 represents volume 5.
	AirVolume5 = "5"
	// AirVolume6 represents volume 6.
	AirVolume6 = "6"
	// AirVolume7 represents volume 7.
	AirVolume7 = "7"
	// AirVolume8 represents volume 8.
	AirVolume8 = "8"
	// AirVolume9 represents volume 9.
	AirVolume9 = "9"
	// AirVolume10 represents volume 10.
	AirVolume10 = "10"
)

// AirDirection represents direction of air.
type AirDirection string

// StringValue converts AirDirection to string.
func (d AirDirection) StringValue() string {
	return string(d)
}

const (
	// AirDirectionAuto represents auto mode.
	AirDirectionAuto AirDirection = ""
)

// Button represents button of air conditioner such as power-off button and power-on button.
type Button string

// StringValue converts Button to string.
func (b Button) StringValue() string {
	return string(b)
}

const (
	// ButtonPowerOn represents power-on button.
	ButtonPowerOn Button = ""
	// ButtonPowerOff represents power-off button.
	ButtonPowerOff Button = "power-off"
)
