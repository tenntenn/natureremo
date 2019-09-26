package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	dialogflow "github.com/leboncoin/dialogflow-go-webhook"
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

func (s *Server) handleBotMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, pass, ok := r.BasicAuth()
	if !(ok && user == "tenntenn" && pass == "tenntenn") {
		code := http.StatusUnauthorized
		http.Error(w, http.StatusText(code), code)
		return
	}

	defer r.Body.Close()
	var dfr dialogflow.Request
	if err := json.NewDecoder(r.Body).Decode(&dfr); err != nil {
		code := http.StatusBadRequest
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	switch dfr.QueryResult.Intent.DisplayName {
	case "Reserve":
		s.reserve(ctx, w, &dfr)
	case "List":
		s.list(ctx, w, &dfr)
	}
}

func (s *Server) reserve(ctx context.Context, w http.ResponseWriter, r *dialogflow.Request) {
	var param DialogflowParam
	if err := json.Unmarshal([]byte(r.QueryResult.Parameters), &param); err != nil {
		code := http.StatusBadRequest
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	sch, err := param.ToSchedule()
	if err != nil {
		code := http.StatusInternalServerError
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	if err := s.scheduler.Register(ctx, sch); err != nil {
		code := http.StatusInternalServerError
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	msg := sch.String() + "にします"
	dff := &dialogflow.Fulfillment{
		FulfillmentMessages: dialogflow.Messages{
			dialogflow.Message{
				Platform:    dialogflow.Line,
				RichMessage: dialogflow.SingleSimpleResponse(msg, msg),
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dff); err != nil {
		log.Println("Error:", err)
	}
}

func (s *Server) list(ctx context.Context, w http.ResponseWriter, r *dialogflow.Request) {
}
