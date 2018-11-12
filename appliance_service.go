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
