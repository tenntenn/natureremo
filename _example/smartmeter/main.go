package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tenntenn/natureremo"
)

type Item struct {
	Value float64
	Time  time.Time
}

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()
	prevNormalCumulative := map[string]Item{}
	prevReverseCumulative := map[string]Item{}
	for {
		as, err := cli.ApplianceService.GetAll(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, a := range as {
			if a.Type == natureremo.ApplianceTypeSmartMeter {
				fmt.Println("===", a.Device.Name, "===")
				if a.SmartMeter != nil {
					vi, _, e := a.SmartMeter.GetMeasuredInstantaneousWatt()
					if e != nil {
						log.Fatal(e)
					}
					fmt.Println("Instantaneous:", vi, "W")
					v, t, e := a.SmartMeter.GetNormalDirectionCumulativeElectricEnergyWattHour()
					if e != nil {
						log.Fatal(e)
					}
					fmt.Println("NormalCumulative:", v/1000, "kWh")
					if prev, ok := prevNormalCumulative[a.ID]; ok && prev.Time.Before(t) {
						diff, err := a.SmartMeter.CalcCumulativeDiff(v, prev.Value)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("NormalCumulativeDiff:", t.Sub(prev.Time), diff/1000, "kWh")
					}
					prevNormalCumulative[a.ID] = Item{v, t}
					v, t, e = a.SmartMeter.GetReverseDirectionCumulativeElectricEnergyWattHour()
					if e != nil {
						log.Fatal(e)
					}
					fmt.Println("ReverseCumulative:", v/1000, "kWh")
					if prev, ok := prevReverseCumulative[a.ID]; ok && prev.Time.Before(t) {
						diff, err := a.SmartMeter.CalcCumulativeDiff(v, prev.Value)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("ReverseCumulativeDiff:", t.Sub(prev.Time), diff/1000, "kWh")
					}
					prevReverseCumulative[a.ID] = Item{v, t}
				}
			}
		}
		d := cli.LastRateLimit.Reset.Sub(time.Now())
		time.Sleep(d / time.Duration(cli.LastRateLimit.Remaining))
	}
}
