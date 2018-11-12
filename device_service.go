package natureremo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// DeviceService gets devices.
type DeviceService interface {
	// GetAll gets devices.
	GetAll(ctx context.Context) ([]*Device, error)
	Update(ctx context.Context, device *Device) (*Device, error)
	Delete(ctx context.Context, device *Device) error
	UpdateTemperatureOffset(ctx context.Context, device *Device) (*Device, error)
	UpdateHumidityOffset(ctx context.Context, device *Device) (*Device, error)
}

type deviceService struct {
	cli *Client
}

// GetAll send a GET request to /1/devices.
func (s *deviceService) GetAll(ctx context.Context) ([]*Device, error) {
	var ds []*Device
	if err := s.cli.get(ctx, "devices", nil, &ds); err != nil {
		return nil, errors.Wrap(err, "GET deviecs failed")
	}
	return ds, nil
}

func (s *deviceService) Update(ctx context.Context, device *Device) (*Device, error) {
	path := fmt.Sprintf("devices/%s", device.ID)

	data := url.Values{}
	data.Set("name", device.Name)

	var d Device
	if err := s.cli.postForm(ctx, path, data, &d); err != nil {
		return nil, errors.Wrapf(err, "POST %s with %#v", path, data)
	}

	return &d, nil
}

func (s *deviceService) Delete(ctx context.Context, device *Device) error {
	path := fmt.Sprintf("devices/%s/delete", device.ID)
	if err := s.cli.post(ctx, path, nil); err != nil {
		return errors.Wrapf(err, "POST %s failed", path)
	}
	return nil
}

func (s *deviceService) UpdateTemperatureOffset(ctx context.Context, device *Device) (*Device, error) {
	path := fmt.Sprintf("devices/%s/temperature_offset", device.ID)

	data := url.Values{}
	data.Set("offset", strconv.FormatInt(device.TemperatureOffset, 10))

	var d Device
	if err := s.cli.postForm(ctx, path, data, &d); err != nil {
		return nil, errors.Wrapf(err, "POST %s with %#v", path, data)
	}

	return &d, nil
}

func (s *deviceService) UpdateHumidityOffset(ctx context.Context, device *Device) (*Device, error) {
	path := fmt.Sprintf("devices/%s/humidity_offset", device.ID)

	data := url.Values{}
	data.Set("offset", strconv.FormatInt(device.HumidityOffset, 10))

	var d Device
	if err := s.cli.postForm(ctx, path, data, &d); err != nil {
		return nil, errors.Wrapf(err, "POST %s with %#v", path, data)
	}

	return &d, nil
}
