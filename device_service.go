package natureremo

import (
	"context"

	"github.com/pkg/errors"
)

type DeviceService interface {
	Devices(ctx context.Context) ([]*Device, error)
}

type deviceService struct {
	cli *Client
}

func (s *deviceService) Devices(ctx context.Context) ([]*Device, error) {
	var ds []*Device
	if err := s.cli.get(ctx, "devices", nil, &ds); err != nil {
		return nil, errors.Wrap(err, "GET deviecs failed")
	}
	return ds, nil
}
