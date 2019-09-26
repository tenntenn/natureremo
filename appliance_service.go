package natureremo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// ApplianceService provides interface of Nature Remo APIs which are related to appliances.
type ApplianceService interface {
	// Detect detects air conditioners by an infrared ray signal.
	Detect(ctx context.Context, ir *IRSignal) ([]*DetectedAirCon, error)
	// GetAll gets all appliances.
	GetAll(ctx context.Context) ([]*Appliance, error)
	// New creates new appliance and links to specified device.
	New(ctx context.Context, device *Device, nickname, image string) (*Appliance, error)
	// NewWithModel creates new appliance with specified model.
	NewWithModel(ctx context.Context, device *Device, nickname, image string, model *ApplianceModel) (*Appliance, error)
	// ReOrder arranges appliances by given orders.
	ReOrder(ctx context.Context, appliances []*Appliance) error
	// Delete deletes specified appliance.
	Delete(ctx context.Context, appliance *Appliance) error
	// Update updates specified appliance.
	Update(ctx context.Context, appliance *Appliance) (*Appliance, error)
	// UpdateAirConSettings updates air conditioner settings of specified appliance.
	UpdateAirConSettings(ctx context.Context, appliance *Appliance, settings *AirConSettings) error
	// SendTVSignal sends TV infrared signal.
	SendTVSignal(ctx context.Context, appliance *Appliance, buttonName string) (*TVState, error)
	// SendLightSignal sends light infrared signal.
	SendLightSignal(ctx context.Context, appliance *Appliance, buttonName string) (*LightState, error)
}

type applianceService struct {
	cli *Client
}

func (s *applianceService) Detect(ctx context.Context, ir *IRSignal) ([]*DetectedAirCon, error) {
	data := url.Values{}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ir); err != nil {
		return nil, errors.Wrapf(err, "cannot encode IRSignal %v", ir)
	}

	data.Set("message", buf.String())
	var aircons []*DetectedAirCon
	if err := s.cli.postForm(ctx, "detectappliance", data, &aircons); err != nil {
		return nil, errors.Wrapf(err, "POST detectappliance failed with %#v", ir)
	}
	return aircons, nil
}

func (s *applianceService) GetAll(ctx context.Context) ([]*Appliance, error) {
	var as []*Appliance
	if err := s.cli.get(ctx, "appliances", nil, &as); err != nil {
		return nil, errors.Wrap(err, "GET appliances failed")
	}
	return as, nil
}

func (s *applianceService) New(ctx context.Context, device *Device, nickname, image string) (*Appliance, error) {
	return s.NewWithModel(ctx, device, nickname, image, nil)
}

func (s *applianceService) NewWithModel(ctx context.Context, device *Device, nickname, image string, model *ApplianceModel) (*Appliance, error) {
	data := url.Values{}
	data.Set("nickname", nickname)
	if model != nil {
		data.Set("model", model.ID)
	}
	data.Set("device", device.ID)
	data.Set("image", image)

	var a Appliance
	if err := s.cli.postForm(ctx, "appliances", data, &a); err != nil {
		return nil, errors.Wrapf(err, "POST appliances failed with %#v", data)
	}
	return &a, nil
}

func (s *applianceService) ReOrder(ctx context.Context, appliances []*Appliance) error {
	ids := make([]string, 0, len(appliances))
	for i := range appliances {
		ids = append(ids, appliances[i].ID)
	}

	data := url.Values{}
	data.Set("appliances", strings.Join(ids, ","))

	if err := s.cli.postForm(ctx, "appliance_orders", data, nil); err != nil {
		return errors.Wrapf(err, "POST appliance_orders failed with %#v", data)
	}
	return nil
}

func (s *applianceService) Delete(ctx context.Context, appliance *Appliance) error {
	path := fmt.Sprintf("appliances/%s/delete", appliance.ID)
	if err := s.cli.post(ctx, path, nil); err != nil {
		return errors.Wrapf(err, "POST %s", path)
	}
	return nil
}

func (s *applianceService) Update(ctx context.Context, appliance *Appliance) (*Appliance, error) {
	path := fmt.Sprintf("appliances/%s", appliance.ID)

	data := url.Values{}
	data.Set("image", appliance.Image)
	data.Set("nickname", appliance.Nickname)

	var a Appliance
	if err := s.cli.postForm(ctx, path, data, &a); err != nil {
		return nil, errors.Wrapf(err, "POST %s with %#v", path, data)
	}
	return &a, nil
}

func (s *applianceService) UpdateAirConSettings(ctx context.Context, appliance *Appliance, settings *AirConSettings) error {
	path := fmt.Sprintf("appliances/%s/aircon_settings", appliance.ID)

	data := url.Values{}
	data.Set("temperature", settings.Temperature)
	data.Set("operation_mode", settings.OperationMode.StringValue())
	data.Set("air_volume", settings.AirVolume.StringValue())
	data.Set("air_direction", settings.AirDirection.StringValue())
	data.Set("button", settings.Button.StringValue())

	if err := s.cli.postForm(ctx, path, data, nil); err != nil {
		return errors.Wrapf(err, "POST %s with %#v", path, data)
	}
	return nil
}

func (s *applianceService) SendTVSignal(ctx context.Context, appliance *Appliance, buttonName string) (*TVState, error) {
	path := fmt.Sprintf("appliances/%s/tv", appliance.ID)

	data := url.Values{}
	data.Set("button", buttonName)

	var status TVState
	if err := s.cli.postForm(ctx, path, data, &status); err != nil {
		return nil, errors.Wrapf(err, "POST %s with %#v", path, data)
	}
	return &status, nil
}

func (s *applianceService) SendLightSignal(ctx context.Context, appliance *Appliance, buttonName string) (*LightState, error) {
	path := fmt.Sprintf("appliances/%s/light", appliance.ID)

	data := url.Values{}
	data.Set("button", buttonName)

	var status LightState
	if err := s.cli.postForm(ctx, path, data, &status); err != nil {
		return nil, errors.Wrapf(err, "POST %s with %#v", path, data)
	}
	return &status, nil
}
