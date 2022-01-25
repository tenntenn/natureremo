package natureremo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// DeviceService provides interface of Nature Remo APIs which are related to devices.
type DeviceService interface {
	// GetAll gets all information of devices which are related with user.
	GetAll(ctx context.Context) ([]*Device, error)
	// Update updates device information which exclude temperature offset and humidity offset.
	Update(ctx context.Context, device *Device) (*Device, error)
	// Delete deletes specified device.
	Delete(ctx context.Context, device *Device) error
	// UpdateTemperatureOffset updates temperature offset of specified device.
	UpdateTemperatureOffset(ctx context.Context, device *Device) (*Device, error)
	// UpdateHumidityOffset updates humidity offset of specified device.
	UpdateHumidityOffset(ctx context.Context, device *Device) (*Device, error)
}

type deviceService struct {
	cli *Client
}

func (s *deviceService) GetAll(ctx context.Context) ([]*Device, error) {
	var ds []*Device
	if err := s.cli.get(ctx, "devices", nil, &ds); err != nil {
		return nil, fmt.Errorf("GET devices failed: %w", err)
	}
	return ds, nil
}

func (s *deviceService) Update(ctx context.Context, device *Device) (*Device, error) {
	path := fmt.Sprintf("devices/%s", device.ID)

	data := url.Values{}
	data.Set("name", device.Name)

	var d Device
	if err := s.cli.postForm(ctx, path, data, &d); err != nil {
		return nil, fmt.Errorf("POST %s with %#v: %w", path, data, err)
	}

	return &d, nil
}

func (s *deviceService) Delete(ctx context.Context, device *Device) error {
	path := fmt.Sprintf("devices/%s/delete", device.ID)
	if err := s.cli.post(ctx, path, nil); err != nil {
		return fmt.Errorf("POST %s failed: %w", path, err)
	}
	return nil
}

func (s *deviceService) UpdateTemperatureOffset(ctx context.Context, device *Device) (*Device, error) {
	path := fmt.Sprintf("devices/%s/temperature_offset", device.ID)

	data := url.Values{}
	data.Set("offset", strconv.FormatInt(device.TemperatureOffset, 10))

	var d Device
	if err := s.cli.postForm(ctx, path, data, &d); err != nil {
		return nil, fmt.Errorf("POST %s with %#v: %w", path, data, err)
	}

	return &d, nil
}

func (s *deviceService) UpdateHumidityOffset(ctx context.Context, device *Device) (*Device, error) {
	path := fmt.Sprintf("devices/%s/humidity_offset", device.ID)

	data := url.Values{}
	data.Set("offset", strconv.FormatInt(device.HumidityOffset, 10))

	var d Device
	if err := s.cli.postForm(ctx, path, data, &d); err != nil {
		return nil, fmt.Errorf("POST %s with %#v: %w", path, data, err)
	}

	return &d, nil
}
