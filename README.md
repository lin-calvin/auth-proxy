# Auth Proxy

A lightweight authentication middleware proxy for web applications. Sits in front of your backend service and provides a login screen with JWT-based authentication.

## Features

- Reverse proxy with authentication layer
- Material Design login page
- JWT token-based authentication
- Configurable via YAML
- OIDC-ready architecture (pluggable auth providers)
- Kubernetes-ready with provided manifests

## Quick Start

```bash
# Build and run
go mod download
go run ./cmd/server

# Access at http://localhost:8080
# Default credentials: admin / admin123
```

## Configuration

Configuration is loaded from `config.yaml` (or path specified by `CONFIG_PATH` env var):

```yaml
server:
  listen: ":8080"

backend:
  url: "http://localhost:3000"

auth:
  jwt_secret: "your-secret-key"
  cookie_name: "auth_token"
  cookie_secure: false
  cookie_max_age: 24h
  token_duration: 24h

users:
  - username: "admin"
    password_hash: "$2a$10$..."
```

### Generating Password Hashes

```bash
go run tools/genhash.go your-password
```

## Endpoints

| Route | Method | Description |
|-------|--------|-------------|
| `/login` | GET | Redirect to login page |
| `/login` | POST | Authenticate (`{username, password}`) |
| `/login-page` | GET | Login form |
| `/logout` | GET | Clear auth cookie |
| `/static/*` | GET | Static assets |
| `/*` | ALL | Protected (proxied to backend) |

## Docker

```bash
docker build -t auth-proxy:latest .
docker run -v $(pwd)/config.yaml:/app/config.yaml -p 8080:8080 auth-proxy:latest
```

## Kubernetes

```bash
kubectl apply -f k8s.yaml
```

Update the `config.yaml` in the ConfigMap with your backend URL and user credentials.

## Architecture

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   Client    │────▶│  Auth Proxy      │────▶│  Backend Service│
│  (Browser)  │     │  (Go)            │     │                 │
└─────────────┘     └──────────────────┘     └─────────────────┘
                          │
                    ┌─────┴─────┐
                    │ config    │
                    │ .yaml     │
                    └───────────┘
```

## Extending for OIDC

Implement the `auth.Provider` interface:

```go
type Provider interface {
    Authenticate(ctx context.Context, username, password string) (*User, error)
}
```

Create `internal/auth/oidc.go` and register in `cmd/server/main.go`.

## License

MIT
