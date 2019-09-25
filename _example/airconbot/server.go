package main

import (
	"net/http"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Server struct {
	initOnce     sync.Once
	router       http.ServeMux
	scheduler    *Scheduler
	bot          *linebot.Client
	webhookToken string
}

func NewServer(sch *Scheduler, bot *linebot.Client, webhookToken string) *Server {
	s := &Server{
		scheduler:    sch,
		bot:          bot,
		webhookToken: webhookToken,
	}
	s.initOnce.Do(s.initHandler)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) initHandler() {
	s.router.HandleFunc("/cron/checkAndRun", s.handleCheckAndRun)
	s.router.HandleFunc("/register", s.handleRegister)
	s.router.HandleFunc("/bot/message", s.handleBotMessage)
}
