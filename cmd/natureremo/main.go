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
		if a.Nickname == "テレビ" {
			for _, s := range a.Signals {
				if s.Name == "日テレ" {
					if err := cli.SignalService.Send(ctx, s); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
