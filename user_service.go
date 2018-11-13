package natureremo

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

// UserService provides interface of Nature Remo APIs which are related to user.
type UserService interface {
	// Me gets own user data.
	Me(ctx context.Context) (*User, error)
	// Update updates user data.
	Update(ctx context.Context, u *User) (*User, error)
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

func (s *userService) Update(ctx context.Context, me *User) (*User, error) {
	data := url.Values{}
	data.Set("nickname", me.Nickname)
	var u User
	if err := s.cli.postForm(ctx, "users/me", data, &u); err != nil {
		return nil, errors.Wrapf(err, "POST users/me failed with %#v", me)
	}
	return &u, nil
}
