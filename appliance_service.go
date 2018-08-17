package natureremo

import "context"

type ApplianceService interface {
	GetAll(ctx context.Context) ([]*Appliance, error)
	GetOrders(ctx context.Context)
	New(ctx context.Context, nickname, model, device, image string) (*Appliance, error)
	Detect(ctx context.Context, is *IRSignal) error
}
