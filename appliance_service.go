package natureremo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

type ApplianceService interface {
	Detect(ctx context.Context, ir *IRSignal) ([]*DetectedAircon, error)
	GetAll(ctx context.Context) ([]*Appliance, error)
	New(ctx context.Context, device *Device, nickname, image string) (*Appliance, error)
	NewWithModel(ctx context.Context, device *Device, nickname, model, image string) (*Appliance, error)
	//GetOrders(ctx context.Context, appliances []*Appliance) ([]*Appliance, error)
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
