package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()
	for {
		ds, err := cli.DeviceService.GetAll(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, d := range ds {
			if len(ds) > 1 {
				fmt.Println("===", d.Name, "===")
			}

			fmt.Println("Temperature:", d.NewestEvents[natureremo.SensorTypeTemperature].Value, "°C")
			fmt.Println("Humidity:", d.NewestEvents[natureremo.SensorTypeHumidity].Value, "%")
			fmt.Println("illumination:", d.NewestEvents[natureremo.SensortypeIllumination].Value)
			fmt.Println("")
		}

		d := cli.LastRateLimit.Reset.Sub(time.Now())
		time.Sleep(d / time.Duration(cli.LastRateLimit.Remaining))
	}
}
