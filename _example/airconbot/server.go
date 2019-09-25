package main

import (
	"net/http"
	"sync"
)

type Server struct {
	initOnce     sync.Once
	router       http.ServeMux
	scheduler    *Scheduler
	webhookToken string
}

func NewServer(sch *Scheduler, webhookToken string) *Server {
	s := &Server{
		scheduler:    sch,
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
}
