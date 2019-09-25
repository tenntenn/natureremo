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
	sch := NewScheduler(ncli)
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_ACCESS_TOKEN"),
	)
	if err != nil {
		panic(err)
	}
	webhookToken := os.Getenv("WEBHOOK_TOKEN")
	server := NewServer(sch, bot, webhookToken)

	http.ListenAndServe(hostport, server)
}
