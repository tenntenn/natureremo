package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dialogflow "github.com/leboncoin/dialogflow-go-webhook"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mercari.io/datastore"
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

func (s *Server) basicAuth(r *http.Request) bool {
	user, pass, ok := r.BasicAuth()
	return ok &&
		s.DialogflowBasicAuth == BasicAuth{User: user, Pass: pass}
}

func (s *Server) handleBotMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !s.basicAuth(r) {
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
	case "Remove":
		s.remove(ctx, w, &dfr)
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

	msg := linebot.NewTextMessage(sch.String() + "にします")
	dff := &dialogflow.Fulfillment{
		FulfillmentMessages: dialogflow.Messages{
			dialogflow.Message{
				Platform: dialogflow.Line,
				RichMessage: dialogflow.PayloadWrapper{Payload: map[string]interface{}{
					"line": msg,
				}},
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dff); err != nil {
		log.Println("Error:", err)
	}
}

func (s *Server) list(ctx context.Context, w http.ResponseWriter, r *dialogflow.Request) {
	schs, err := s.scheduler.GetAll(ctx)
	if err != nil && err != datastore.ErrNoSuchEntity {
		code := http.StatusInternalServerError
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	if len(schs) == 0 {
		dff := &dialogflow.Fulfillment{
			FulfillmentMessages: dialogflow.Messages{
				dialogflow.Message{
					Platform: dialogflow.Line,
					RichMessage: dialogflow.PayloadWrapper{Payload: map[string]interface{}{
						"line": linebot.NewTextMessage("予約はありません"),
					}},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dff); err != nil {
			log.Println("Error:", err)
		}
		return
	}

	var cols []*linebot.CarouselColumn
	done := map[string]bool{}
	for i := range schs {
		key := schs[i].String()
		date, aircon, mode, button := schs[i].Format()
		body := fmt.Sprintf("%sの%sを%s", aircon, mode, button)
		if done[key] {
			continue
		}
		action := linebot.NewMessageAction("削除", fmt.Sprintf("%sを削除", key))
		cols = append(cols, linebot.NewCarouselColumn("", date, body, action))
		done[key] = true
	}

	template := linebot.NewCarouselTemplate(cols...)
	dff := &dialogflow.Fulfillment{
		FulfillmentMessages: dialogflow.Messages{
			dialogflow.Message{
				Platform: dialogflow.Line,
				RichMessage: dialogflow.PayloadWrapper{Payload: map[string]interface{}{
					"line": linebot.NewTemplateMessage("予約一覧", template),
				}},
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dff); err != nil {
		log.Println("Error:", err)
	}
}

func (s *Server) remove(ctx context.Context, w http.ResponseWriter, r *dialogflow.Request) {
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

	if err := s.scheduler.Delete(ctx, sch); err != nil {
		code := http.StatusInternalServerError
		log.Println("Error:", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	msg := linebot.NewTextMessage(sch.String() + "を削除しました")
	dff := &dialogflow.Fulfillment{
		FulfillmentMessages: dialogflow.Messages{
			dialogflow.Message{
				Platform: dialogflow.Line,
				RichMessage: dialogflow.PayloadWrapper{Payload: map[string]interface{}{
					"line": msg,
				}},
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dff); err != nil {
		log.Println("Error:", err)
	}
}
