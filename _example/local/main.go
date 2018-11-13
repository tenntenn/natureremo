package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tenntenn/natureremo"
)

func main() {
	cli := natureremo.NewLocalClient(os.Args[1])
	ctx := context.Background()

	ir, err := cli.Fetch(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ir)

	if err := cli.Emit(ctx, ir); err != nil {
		log.Fatal(err)
	}
}
