package natureremo

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type SmartMeter struct {
	Properties []EchonetliteProperty `json:"echonetlite_properties"`
}

type EchonetliteProperty struct {
	Name      string         `json:"name"`
	Epc       EchonetliteEPC `json:"epc"`
	Value     string         `json:"val"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type EchonetliteEPC int

const (
	EPCNormalDirectionCumulativeElectricEnergy  EchonetliteEPC = 0xE0
	EPCReverseDirectionCumulativeElectricEnergy EchonetliteEPC = 0xE3
	EPCCoefficient                              EchonetliteEPC = 0xD3
	EPCCumulativeElectricEnergyUnit             EchonetliteEPC = 0xE1
	EPCCumulativeElectricEnergyEffectiveDigits  EchonetliteEPC = 0xD7
	EPCMeasuredInstantaneous                    EchonetliteEPC = 0xE7
)

// Find specified EPC object.
// It returns reference to the object if found. Otherwise, returns nil.
func (s *SmartMeter) Find(epc EchonetliteEPC) *EchonetliteProperty {
	for _, p := range s.Properties {
		if p.Epc == epc {
			return &p
		}
	}
	return nil
}

// Get MeasuredInstantaneous in Watt(W).
func (s *SmartMeter) GetMeasuredInstantaneousWatt() (int64, time.Time, error) {
	p := s.Find(EPCMeasuredInstantaneous)
	if p == nil {
		return 0, time.Time{}, fmt.Errorf("MeasuredInstantaneous property not found")
	}
	v, err := strconv.ParseInt(p.Value, 10, 64)
	if err != nil {
		return 0, time.Time{}, err
	}
	return v, p.UpdatedAt, nil
}

// Get NormalDirectionCumulativeElectricEnergy in watt hour(Wh).
func (s *SmartMeter) GetNormalDirectionCumulativeElectricEnergyWattHour() (float64, time.Time, error) {
	coefficient, err := s.getCumulativeElectricEnergyCoefficientWattHour()
	if err != nil {
		return 0, time.Time{}, err
	}
	p := s.Find(EPCNormalDirectionCumulativeElectricEnergy)
	if p == nil {
		return 0, time.Time{}, fmt.Errorf("NormalDirectionCumulativeElectricEnergy property not found")
	}
	v, err := strconv.ParseUint(p.Value, 10, 64)
	if err != nil {
		return 0, time.Time{}, err
	}
	return float64(v) * coefficient, p.UpdatedAt, nil
}

// Get ReverseDirectionCumulativeElectricEnergy in watt hour(Wh).
func (s *SmartMeter) GetReverseDirectionCumulativeElectricEnergyWattHour() (float64, time.Time, error) {
	coefficient, err := s.getCumulativeElectricEnergyCoefficientWattHour()
	if err != nil {
		return 0, time.Time{}, err
	}
	p := s.Find(EPCReverseDirectionCumulativeElectricEnergy)
	if p == nil {
		return 0, time.Time{}, fmt.Errorf("ReverseDirectionCumulativeElectricEnergy property not found")
	}
	v, err := strconv.ParseUint(p.Value, 10, 64)
	if err != nil {
		return 0, time.Time{}, err
	}
	return float64(v) * coefficient, p.UpdatedAt, nil
}

// Calculate diff of two cumulative value considering over flow.
func (s *SmartMeter) CalcCumulativeDiff(nextWattHour float64, prevWattHour float64) (float64, error) {
	coefficient, err := s.getCumulativeElectricEnergyCoefficientWattHour()
	if err != nil {
		return 0, err
	}
	effectiveDigitsProp := s.Find(EPCCumulativeElectricEnergyEffectiveDigits)
	if effectiveDigitsProp == nil {
		return 0, fmt.Errorf("CumulativeElectricEnergyEffectiveDigits property not found")
	}
	effectiveDigits, err := strconv.ParseUint(effectiveDigitsProp.Value, 10, 8)
	if err != nil {
		return 0, err
	}
	upperBound := math.Pow10(int(effectiveDigits+1)) * coefficient
	if nextWattHour >= prevWattHour {
		return nextWattHour - prevWattHour, nil
	} else {
		return upperBound - prevWattHour + nextWattHour, nil
	}
}

// Get coefficient to convert cumulative electric energy unit to watt hour
func (s *SmartMeter) getCumulativeElectricEnergyCoefficientWattHour() (value float64, err error) {
	unitProp := s.Find(EPCCumulativeElectricEnergyUnit)
	var unitCode uint64
	if unitProp == nil {
		unitCode = 0x00 // assume unit is kWh
	} else {
		unitCode, err = strconv.ParseUint(unitProp.Value, 10, 8)
		if err != nil {
			return 0.0, err
		}
	}
	var unitCoefficient float64
	if unitCode < 0x0A {
		unitCoefficient = math.Pow10(0x03 - int(unitCode))
	} else {
		unitCoefficient = math.Pow10(int(unitCode) - 0x0A + 4)
	}
	coefficientProp := s.Find(EPCCoefficient)
	var coefficient uint64
	if coefficientProp == nil {
		coefficient = 1
	} else {
		coefficient, err = strconv.ParseUint(coefficientProp.Value, 10, 64)
		if err != nil {
			return 0.0, err
		}
	}
	return unitCoefficient * float64(coefficient), nil
}
