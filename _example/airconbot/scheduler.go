package main

import (
	"context"

	"github.com/tenntenn/natureremo"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
	"go.uber.org/multierr"
)

type Schedule struct {
	ScheduledAt int64             `json:"scheduledAt" datastore:"scheduled_at"`
	ApplianceID string            `json:"applianceId" datastore:"appliance_id"`
	Button      natureremo.Button `json:"button" datastore:"button"`
}

type Scheduler struct {
	ncli *natureremo.Client
}

func NewScheduler(ncli *natureremo.Client) *Scheduler {
	return &Scheduler{
		ncli: ncli,
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

		err := s.ncli.ApplianceService.UpdateAirConSettings(ctx, a, &settings)
		if err != nil {
			return err
		}
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

	return nil
}
