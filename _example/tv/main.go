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

	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	const (
		applianceTV   = "テレビ" // means TV in Japan
		signalChannel = "日テレ" // is one of channels in Japan
	)

	var tv *natureremo.Appliance
	for _, a := range as {
		if a.Nickname == applianceTV {
			tv = a
			break
		}
	}

	if tv == nil {
		log.Fatal("TV not found")
	}

	for _, s := range tv.Signals {
		if s.Name == signalChannel {
			cli.SignalService.Send(ctx, s)
			break
		}
	}
}
