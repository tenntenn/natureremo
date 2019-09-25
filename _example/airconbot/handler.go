package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var jstTz = time.FixedZone("JST", 9*60*60)

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

func (s *Server) handleBotMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := s.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			return
		}
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			continue
		}

		switch m := event.Message.(type) {
		case *linebot.TextMessage:
			switch {
			case strings.HasPrefix(m.Text, "エアコンオン"):
				splited := strings.Split(m.Text, " ")
				if len(splited) >= 2 {
					date := strings.Join(splited[1:], " ")
					tm, err := time.ParseInLocation("2006/01/02 15:04", date, jstTz)
					if err != nil {
						code := http.StatusInternalServerError
						log.Println("Error:", err)
						http.Error(w, http.StatusText(code)+err.Error(), code)
						return
					}
					sch := &Schedule{
						ScheduledAt: tm.Unix(),
						ApplianceID: "",
					}
					if err := s.scheduler.Register(ctx, sch); err != nil {
						code := http.StatusInternalServerError
						log.Println("Error:", err)
						http.Error(w, http.StatusText(code)+err.Error(), code)
						return
					}

					msg := linebot.NewTextMessage(date + "でエアコンを予約しました")
					if _, err = s.bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
						code := http.StatusInternalServerError
						log.Println("Error:", err)
						http.Error(w, http.StatusText(code)+err.Error(), code)
						return
					}
				}
			}
		}
	}
}
