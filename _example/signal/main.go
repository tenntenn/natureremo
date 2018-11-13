package main

import (
	"context"
	"log"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()

	applianceName := os.Args[2]
	signalName := os.Args[3]

	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var target *natureremo.Appliance
	for _, a := range as {
		if a.Nickname == applianceName {
			target = a
			break
		}
	}

	if target == nil {
		log.Fatalf("%s not found", applianceName)
	}

	for _, s := range target.Signals {
		if s.Name == signalName {
			cli.SignalService.Send(ctx, s)
			break
		}
	}
}
