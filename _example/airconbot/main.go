package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tenntenn/natureremo"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	hostport := net.JoinHostPort("", port)

	ncli := natureremo.NewClient(os.Getenv("NATUREREMO_TOKEN"))
	bot, err := linebot.New(
		os.Getenv("LINE_SECRET"),
		os.Getenv("LINE_TOKEN"),
	)
	if err != nil {
		return err
	}
	sch := NewScheduler(ncli, bot)
	server := NewServer(sch)

	return http.ListenAndServe(hostport, server)
}
