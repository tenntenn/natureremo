package natureremo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// SignalService provides interface of Nature Remo APIs which are related to signals.
type SignalService interface {
	// GetAll gets all signals which related to specified appliance.
	GetAll(ctx context.Context, appliance *Appliance) ([]*Signal, error)
	// New creates new signal and links to specified appliance.
	New(ctx context.Context, appliance *Appliance, ir *IRSignal, name, image string) (*Signal, error)
	// ReOrder arranges signals by given orders.
	ReOrder(ctx context.Context, appliance *Appliance, signals []*Signal) error
	// Update updates specified signal.
	Update(ctx context.Context, signal *Signal) (*Signal, error)
	// Delete deletes specified signal.
	Delete(ctx context.Context, signal *Signal) error
	// Send sends specified signal.
	Send(ctx context.Context, signal *Signal) error
}

type signalService struct {
	cli *Client
}

func (s *signalService) GetAll(ctx context.Context, appliance *Appliance) ([]*Signal, error) {
	path := fmt.Sprintf("appliances/%s/signals", appliance.ID)
	var ss []*Signal
	if err := s.cli.get(ctx, path, nil, &ss); err != nil {
		return nil, fmt.Errorf("GET %s failed: %w", path, err)
	}
	return ss, nil
}

func (s *signalService) New(ctx context.Context, appliance *Appliance, ir *IRSignal, name, image string) (*Signal, error) {
	path := fmt.Sprintf("appliances/%s/signals", appliance.ID)
	data := url.Values{}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ir); err != nil {
		return nil, fmt.Errorf("cannot encode IRSignal %v: %w", ir, err)
	}
	data.Set("message", buf.String())
	data.Set("name", name)
	data.Set("image", image)

	var sig Signal
	if err := s.cli.postForm(ctx, path, data, &sig); err != nil {
		return nil, fmt.Errorf("POST %s with %#v: %w", path, data, err)
	}
	return &sig, nil
}

func (s *signalService) ReOrder(ctx context.Context, appliance *Appliance, signals []*Signal) error {
	path := fmt.Sprintf("appliances/%s/signal_orders", appliance.ID)

	ids := make([]string, 0, len(signals))
	for i := range signals {
		ids = append(ids, signals[i].ID)
	}

	data := url.Values{}
	data.Set("signals", strings.Join(ids, ","))

	if err := s.cli.postForm(ctx, path, data, nil); err != nil {
		return fmt.Errorf("POST %s failed with %#v: %w", path, data, err)
	}

	return nil
}

func (s *signalService) Update(ctx context.Context, signal *Signal) (*Signal, error) {
	path := fmt.Sprintf("signals/%s", signal.ID)

	data := url.Values{}
	data.Set("name", signal.Name)
	data.Set("image", signal.Image)

	var sig Signal
	if err := s.cli.postForm(ctx, path, data, &sig); err != nil {
		return nil, fmt.Errorf("POST %s failed with %#v: %w", path, signal, err)
	}

	return &sig, nil
}

func (s *signalService) Delete(ctx context.Context, signal *Signal) error {
	path := fmt.Sprintf("signals/%s/delete", signal.ID)
	if err := s.cli.post(ctx, path, nil); err != nil {
		return fmt.Errorf("POST %s failed: %w", path, err)
	}
	return nil
}

func (s *signalService) Send(ctx context.Context, signal *Signal) error {
	path := fmt.Sprintf("signals/%s/send", signal.ID)
	if err := s.cli.post(ctx, path, nil); err != nil {
		return fmt.Errorf("POST %s failed: %w", path, err)
	}
	return nil
}
