package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ds, err := cli.DeviceService.Devices(context.Background())
	if err != nil {
		panic(err)
	}
	for i := range ds {
		fmt.Println(ds[i])
	}
}
