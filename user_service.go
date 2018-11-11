package natureremo

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

// UserService is get or update own user data.
type UserService interface {
	// Me gets own user data.
	Me(ctx context.Context) (*User, error)
	// UpdateNickname updates nickname.
	UpdateNickname(ctx context.Context, nickname string) error
}

type userService struct {
	cli *Client
}

// Me sends GET request to /1/users/me.
func (s *userService) Me(ctx context.Context) (*User, error) {
	var u User
	if err := s.cli.get(ctx, "users/me", nil, &u); err != nil {
		return nil, errors.Wrap(err, "GET users/me failed")
	}
	return &u, nil
}

// UpdateNickname sends POST request to /1/users/me.
func (s *userService) UpdateNickname(ctx context.Context, nickname string) error {
	data := url.Values{}
	data.Set("nickname", nickname)
	if err := s.cli.postForm(ctx, "users/me", data, nil); err != nil {
		return errors.Wrapf(err, "POST users/me failed with %s", nickname)
	}
	return nil
}
