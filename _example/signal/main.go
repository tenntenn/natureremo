package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	ctx := context.Background()

	applianceName := os.Args[2]

	a, err := getAppliance(ctx, cli, applianceName)
	if err != nil {
		log.Fatal(err)
	}

	signalName := os.Args[3]
	s := getSignal(a.Signals, signalName)
	if s == nil {
		log.Fatal("signal not found")
	}

	if err := cli.SignalService.Send(ctx, s); err != nil {
		log.Fatal(err)
	}
}

func getAppliance(ctx context.Context, cli *natureremo.Client, name string) (*natureremo.Appliance, error) {
	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		if a.Nickname == name {
			return a, nil
		}
	}

	return nil, errors.New("appliance not found")
}

func getSignal(ss []*natureremo.Signal, name string) *natureremo.Signal {
	for _, s := range ss {
		if s.Name == name {
			return s
		}
	}
	return nil
}
