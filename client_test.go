package natureremo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tenntenn/natureremo"
)

func TestNewClientInvalidToken(t *testing.T) {
	cli := natureremo.NewClient("--invalid--")
	ctx := context.Background()
	_, err := cli.ApplianceService.GetAll(ctx)
	var aerr *natureremo.APIError
	if !errors.As(err, &aerr) ||
		aerr.Code != 401001 {
		t.Fatal("unexpected error", err)
	}
}
