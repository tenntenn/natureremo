package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewClient(os.Args[1])
	vs, err := cli.ApplianceService.GetAll(context.Background())
	if err != nil {
		panic(err)
	}
	json.NewEncoder(os.Stdout).Encode(vs)
}
