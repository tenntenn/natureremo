package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tenntenn/natureremo"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
	"go.uber.org/multierr"
)

type Schedule struct {
	ScheduledAt   int64                    `datastore:"scheduled_at"`
	ApplianceName string                   `datastore:"appliance_name"`
	ApplianceID   string                   `datastore:"appliance_id"`
	Button        natureremo.Button        `datastore:"button"`
	Mode          natureremo.OperationMode `datastore:"mode"`
}

func (s *Schedule) String() string {
	mode := "エアコン"
	switch s.Mode {
	case natureremo.OperationModeWarm:
		mode = "暖房"
	case natureremo.OperationModeCool:
		mode = "冷房"
	}

	button := "ON"
	if s.Button == natureremo.ButtonPowerOff {
		button = "OFF"
	}

	tm := time.Unix(s.ScheduledAt, 0).In(jstTz)
	date := tm.Format("2006/01/02 15:04")
	return fmt.Sprintf("%sに%sの%sを%s", date, s.ApplianceName, mode, button)
}

type Scheduler struct {
	ncli *natureremo.Client
	bot  *linebot.Client
}

func NewScheduler(ncli *natureremo.Client, bot *linebot.Client) *Scheduler {
	return &Scheduler{
		ncli: ncli,
		bot:  bot,
	}
}

func (s *Scheduler) RunAll(ctx context.Context, now int64) error {
	client, err := clouddatastore.FromContext(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	var schedules []*Schedule
	q := client.NewQuery("Schedule").Filter("scheduled_at <=", now)
	keys, err := client.GetAll(ctx, q, &schedules)
	if err != nil {
		return err
	}

	var (
		errs error
		done []datastore.Key
	)
	for i := range schedules {
		if err := s.Run(ctx, schedules[i]); err != nil {
			errs = multierr.Append(errs, err)
		} else {
			done = append(done, keys[i])
		}
	}

	if err := client.DeleteMulti(ctx, done); err != nil {
		return err
	}

	if errs != nil {
		return errs
	}

	if len(schedules) > 0 {
		strs := make([]string, len(schedules))
		for i := range schedules {
			strs[i] = schedules[i].String()
		}
		msg := linebot.NewTextMessage(strings.Join(strs, "、") + "にしました")
		if _, err := s.bot.BroadcastMessage(linebot.SendingMessage(msg)).WithContext(ctx).Do(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) Run(ctx context.Context, sch *Schedule) error {

	as, err := s.ncli.ApplianceService.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, a := range as {
		if a.ID != sch.ApplianceID || a.AirConSettings == nil {
			continue
		}

		settings := *(a.AirConSettings)
		settings.Button = sch.Button

		if sch.Mode != "" {
			settings.OperationMode = sch.Mode
		}

		err := s.ncli.ApplianceService.UpdateAirConSettings(ctx, a, &settings)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (s *Scheduler) Register(ctx context.Context, sch *Schedule) error {
	client, err := clouddatastore.FromContext(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	key := client.IncompleteKey("Schedule", nil)
	if _, err := client.Put(ctx, key, sch); err != nil {
		return err
	}

	log.Println(sch)

	return nil
}
