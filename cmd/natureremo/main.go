package main

import (
	"context"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()

	//ds, err := cli.DeviceService.GetAll(ctx)
	//if err != nil {
	//	panic(err)
	//}

	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	for _, a := range as {
		if a.AirConSettings != nil {
			s := *(a.AirConSettings)
			s.OperationMode = natureremo.OperationModeWarm
			cli.ApplianceService.UpdateAirConSettings(ctx, a, &s)
			break
		}
	}
}
