package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/tenntenn/natureremo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type AirConCalendar struct {
	NatureRemo *natureremo.Client
	Secret     string
}

type Request struct {
	Secret        string `json:"secret"`
	OperationMode string `json:"mode"`
	OnOff         string `json:"onoff"`
}

func (c *AirConCalendar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	c.NatureRemo.HTTPClient = urlfetch.Client(ctx)

	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	defer r.Body.Close()

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code := http.StatusInternalServerError
		log.Errorf(ctx, "cannot decode request body: %s", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	if req.Secret != c.Secret {
		code := http.StatusForbidden
		log.Errorf(ctx, "secret does not much")
		http.Error(w, http.StatusText(code), code)
		return
	}

	as, err := c.NatureRemo.ApplianceService.GetAll(ctx)
	if err != nil {
		code := http.StatusInternalServerError
		log.Errorf(ctx, "%s", err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	for _, a := range as {
		if a.AirConSettings != nil {
			settings := *(a.AirConSettings)

			switch req.OperationMode {
			case "暖房":
				settings.OperationMode = natureremo.OperationModeWarm
			case "冷房":
				settings.OperationMode = natureremo.OperationModeCool
			default:
				settings.OperationMode = natureremo.OperationModeAuto
			}

			switch req.OnOff {
			case "オフ", "off", "OFF", "Off":
				settings.Button = natureremo.ButtonPowerOff
			default:
				settings.Button = natureremo.ButtonPowerOn
				settings.Temperature = ""
			}

			err := c.NatureRemo.ApplianceService.UpdateAirConSettings(ctx, a, &settings)
			if err != nil {
				code := http.StatusInternalServerError
				log.Errorf(ctx, "%s", err)
				http.Error(w, http.StatusText(code), code)
				return
			}
			log.Infof(ctx, "%#v", req)
			break
		}
	}
}

func init() {
	c := &AirConCalendar{
		NatureRemo: natureremo.NewClient(os.Getenv("NATUREREMO_TOKEN")),
		Secret:     os.Getenv("WEBHOOK_SECRET"),
	}
	http.Handle("/", c)
}
