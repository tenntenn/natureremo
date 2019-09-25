package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleCheckAndRun(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-AppEngine-Cron") == "" {
		code := http.StatusForbidden
		http.Error(w, http.StatusText(code), code)
		return
	}

	ctx := r.Context()
	if err := s.scheduler.RunAll(ctx, time.Now().Unix()); err != nil {
		code := http.StatusInternalServerError
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Header.Get("X-AIRCONBOT-TOKEN") != s.webhookToken {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
	*/
	defer r.Body.Close()

	var sch Schedule
	if err := json.NewDecoder(r.Body).Decode(&sch); err != nil {
		code := http.StatusBadRequest
		http.Error(w, http.StatusText(code), code)
		return
	}

	ctx := r.Context()
	if err := s.scheduler.Register(ctx, &sch); err != nil {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}
}
