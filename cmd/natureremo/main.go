package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()

	ds, err := cli.DeviceService.Devices(ctx)
	if err != nil {
		panic(err)
	}

	a, err := cli.ApplianceService.New(ctx, ds[0], "test", "ico_ac_1")
	if err != nil {
		panic(err)
	}
	json.NewEncoder(os.Stdout).Encode(a)
}
