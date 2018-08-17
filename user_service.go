package natureremo

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

type UserService interface {
	Me(ctx context.Context) (*User, error)
	UpdateNickname(ctx context.Context, nickname string) error
}

type userService struct {
	cli *Client
}

func (s *userService) Me(ctx context.Context) (*User, error) {
	var u User
	if err := s.cli.get(ctx, "users/me", nil, &u); err != nil {
		return nil, errors.Wrap(err, "GET users/me failed")
	}
	return &u, nil
}

func (s *userService) UpdateNickname(ctx context.Context, nickname string) error {
	data := url.Values{}
	data.Set("nickname", nickname)
	if err := s.cli.postForm(ctx, "users/me", data, nil); err != nil {
		return errors.Wrapf(err, "POST users/me failed with %s", nickname)
	}
	return nil
}
