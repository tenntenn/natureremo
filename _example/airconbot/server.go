package main

import (
	"net/http"
	"os"
	"sync"
)

type BasicAuth struct {
	User string
	Pass string
}

type Server struct {
	initOnce            sync.Once
	router              http.ServeMux
	scheduler           *Scheduler
	DialogflowBasicAuth BasicAuth
}

func NewServer(sch *Scheduler) *Server {
	s := &Server{
		scheduler: sch,
		DialogflowBasicAuth: BasicAuth{
			User: os.Getenv("DIALOGFLOW_BASIC_AUTH_USER"),
			Pass: os.Getenv("DIALOGFLOW_BASIC_AUTH_PASS"),
		},
	}
	s.initOnce.Do(s.initHandler)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) initHandler() {
	s.router.HandleFunc("/cron/checkAndRun", s.handleCheckAndRun)
	s.router.HandleFunc("/bot/message", s.handleBotMessage)
}
