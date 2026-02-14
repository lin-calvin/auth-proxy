package auth

import (
	"context"
	"errors"

	"auth-proxy/internal/config"

	"golang.org/x/crypto/bcrypt"
)

type StaticProvider struct {
	users map[string]string
}

func NewStaticProvider(userList []config.User) *StaticProvider {
	users := make(map[string]string)
	for _, u := range userList {
		users[u.Username] = u.PasswordHash
	}
	return &StaticProvider{users: users}
}

func (p *StaticProvider) Authenticate(ctx context.Context, username, password string) (*User, error) {
	hash, ok := p.users[username]
	if !ok {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &User{
		Username: username,
		Roles:    []string{"user"},
		Claims:   make(map[string]interface{}),
	}, nil
}
