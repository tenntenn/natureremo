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

type ApplianceService interface {
	Detect(ctx context.Context, ir *IRSignal) ([]*DetectedAircon, error)
	GetAll(ctx context.Context) ([]*Appliance, error)
	New(ctx context.Context, device *Device, nickname, image string) (*Appliance, error)
	NewWithModel(ctx context.Context, device *Device, nickname, model, image string) (*Appliance, error)
	ReOrder(ctx context.Context, appliances []*Appliance) error
	Delete(ctx context.Context, appliance *Appliance) error
	Update(ctx context.Context, appliance *Appliance) (*Appliance, error)
	UpdateAirConSettings(ctx context.Context, appliance *Appliance, settings *AirConParams) error
}

type applianceService struct {
	cli *Client
}

func (s *applianceService) Detect(ctx context.Context, ir *IRSignal) ([]*DetectedAircon, error) {
	data := url.Values{}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ir); err != nil {
		return nil, errors.Wrapf(err, "cannot encode IRSignal %v", ir)
	}

	data.Set("message", buf.String())
	var aircons []*DetectedAircon
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
	return s.NewWithModel(ctx, device, nickname, image, "")
}

func (s *applianceService) NewWithModel(ctx context.Context, device *Device, nickname, image, model string) (*Appliance, error) {
	data := url.Values{}
	data.Set("nickname", nickname)
	if model != "" {
		data.Set("model", model)
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

func (s *applianceService) UpdateAirConSettings(ctx context.Context, appliance *Appliance, settings *AirConParams) error {
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
