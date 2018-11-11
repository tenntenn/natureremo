package natureremo

import (
	"context"

	"github.com/pkg/errors"
)

// DeviceService gets devices.
type DeviceService interface {
	// Devices gets devices.
	Devices(ctx context.Context) ([]*Device, error)
}

type deviceService struct {
	cli *Client
}

// Devices send a GET request to /1/devices.
func (s *deviceService) Devices(ctx context.Context) ([]*Device, error) {
	var ds []*Device
	if err := s.cli.get(ctx, "devices", nil, &ds); err != nil {
		return nil, errors.Wrap(err, "GET deviecs failed")
	}
	return ds, nil
}
