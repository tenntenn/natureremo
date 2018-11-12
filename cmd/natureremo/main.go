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

	ds, err := cli.DeviceService.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	ds[0].HumidityOffset = 0

	d, err := cli.DeviceService.UpdateHumidityOffset(ctx, ds[0])
	if err != nil {
		panic(err)
	}
	json.NewEncoder(os.Stdout).Encode(d)
}
