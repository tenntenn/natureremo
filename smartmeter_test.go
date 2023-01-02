package natureremo_test

import (
	"encoding/json"
	"testing"

	"github.com/tenntenn/natureremo"
)

func TestSmartMeter(t *testing.T) {
	jsonString := `{
		"echonetlite_properties": [
		  {
			"name": "coefficient",
			"epc": 211,
			"val": "1",
			"updated_at": "2020-04-27T15:24:03Z"
		  },
		  {
			"name": "cumulative_electric_energy_effective_digits",
			"epc": 215,
			"val": "6",
			"updated_at": "2020-04-27T15:24:03Z"
		  },
		  {
			"name": "normal_direction_cumulative_electric_energy",
			"epc": 224,
			"val": "5167",
			"updated_at": "2020-04-27T15:24:03Z"
		  },
		  {
			"name": "cumulative_electric_energy_unit",
			"epc": 225,
			"val": "1",
			"updated_at": "2020-04-27T15:24:03Z"
		  },
		  {
			"name": "reverse_direction_cumulative_electric_energy",
			"epc": 227,
			"val": "3606",
			"updated_at": "2020-04-27T15:24:03Z"
		  },
		  {
			"name": "measured_instantaneous",
			"epc": 231,
			"val": "360",
			"updated_at": "2020-04-27T15:24:03Z"
		  }
		]
	}`
	s := natureremo.SmartMeter{}
	if err := json.Unmarshal([]byte(jsonString), &s); err != nil {
		t.Fatal(err)
	}
	instant, _, err := s.GetMeasuredInstantaneousWatt()
	if err != nil {
		t.Fatal(err)
	}
	if instant != 360 {
		t.Fatalf("instant value expects 360, but %v", instant)
	}
	normal, _, err := s.GetNormalDirectionCumulativeElectricEnergyWattHour()
	if err != nil {
		t.Fatal(err)
	}
	if normal != 516700.0 {
		t.Fatalf("normal cumulative value expects 516700, but %v", normal)
	}
	reverse, _, err := s.GetReverseDirectionCumulativeElectricEnergyWattHour()
	if err != nil {
		t.Fatal(err)
	}
	if reverse != 360600.0 {
		t.Fatalf("reverse cumulative value expects 360600, but %v", reverse)
	}
}
