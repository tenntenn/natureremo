package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tenntenn/natureremo"
)

var dfAircon = map[string]string{}

func init() {
	for i := 1; ; i++ {
		name := os.Getenv(fmt.Sprintf("AIRCON%d_NAME", i))
		id := os.Getenv(fmt.Sprintf("AIRCON%d_ID", i))
		if name == "" || id == "" {
			if i == 1 {
				panic("AIRCON1_NAME and AIRCON1_ID must be set")
			}
			break
		}
		dfAircon[name] = id
	}
}

var dfButton = map[string]natureremo.Button{
	"ON":  natureremo.ButtonPowerOn,
	"OFF": natureremo.ButtonPowerOff,
}

var dfMode = map[string]natureremo.OperationMode{
	"冷房": natureremo.OperationModeCool,
	"暖房": natureremo.OperationModeWarm,
	"除湿": natureremo.OperationModeDry,
}

type DialogflowParam struct {
	AirCon string `json:"aircon"`
	Mode   string `json:"mode"`
	Button string `json:"button"`
	Date   string `json:"date"`
	Time   string `json:"time"`
}

var jstTz = time.FixedZone("JST", 9*60*60)

func (p *DialogflowParam) ScheduledAt() (time.Time, error) {
	datetime := time.Now().In(jstTz).Format("2006-01-02") + p.Time[len("2006-01-02"):]
	if p.Date != "" {
		datetime = p.Date[:len("2006-01-02")] + p.Time[len("2006-01-02"):]
	}
	return time.ParseInLocation(time.RFC3339, datetime, jstTz)
}

func (p *DialogflowParam) ToSchedule() (*Schedule, error) {
	var sch Schedule

	scheduledAt, err := p.ScheduledAt()
	if err != nil {
		return nil, err
	}
	sch.ScheduledAt = scheduledAt.Unix()

	sch.ApplianceID = os.Getenv("AIRCON1_ID")
	sch.ApplianceName = os.Getenv("AIRCON1_NAME")
	if p.AirCon != "" {
		sch.ApplianceID = dfAircon[p.AirCon]
		sch.ApplianceName = p.AirCon
	}

	sch.Button = natureremo.ButtonPowerOn
	if p.Button != "" {
		sch.Button = dfButton[p.Button]
	}

	return &sch, nil
}
