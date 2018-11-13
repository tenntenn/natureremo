package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/tenntenn/natureremo"
)

func getAddr() (string, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return "", err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = resolver.Browse(ctx, "_remo._tcp", "local.", entries)
	if err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		if err != nil {
			return "", err
		}
	case entry := <-entries:
		ip := entry.AddrIPv4[0]
		addr := net.JoinHostPort(ip.String(), strconv.Itoa(entry.Port))
		return addr, nil
	}

	return "", errors.New("cannot find Nature Remo")
}

func main() {
	addr, err := getAddr()
	fmt.Printf("access to Nature Remo (%s)\n", addr)
	if err != nil {
		log.Fatal(err)
	}

	cli := natureremo.NewLocalClient(addr)
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
