package auth

import "context"

type User struct {
	Username string
	Roles    []string
	Claims   map[string]interface{}
}

type Provider interface {
	Authenticate(ctx context.Context, username, password string) (*User, error)
}
