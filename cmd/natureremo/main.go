package main

import (
	"context"
	"fmt"
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

	a, err := cli.ApplianceService.New(ctx, ds[0], "test", "ico_aircon")
	if err != nil {
		panic(err)
	}

	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(as))
	//json.NewEncoder(os.Stdout).Encode(as)

	err = cli.ApplianceService.Delete(ctx, a)
	if err != nil {
		panic(err)
	}

	as, err = cli.ApplianceService.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(as))
}
