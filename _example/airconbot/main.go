package main

import (
	"net"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tenntenn/natureremo"
)

func main() {
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
		panic(err)
	}
	sch := NewScheduler(ncli, bot)
	server := NewServer(sch)

	http.ListenAndServe(hostport, server)
}
