package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Backend BackendConfig `yaml:"backend"`
	Auth    AuthConfig    `yaml:"auth"`
	Users   []User        `yaml:"users"`
}

type ServerConfig struct {
	Listen string `yaml:"listen"`
}

type BackendConfig struct {
	URL string `yaml:"url"`
}

type AuthConfig struct {
	JWTSecret     string        `yaml:"jwt_secret"`
	CookieName    string        `yaml:"cookie_name"`
	CookieSecure  bool          `yaml:"cookie_secure"`
	CookieMaxAge  time.Duration `yaml:"cookie_max_age"`
	TokenDuration time.Duration `yaml:"token_duration"`
}

type User struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Listen == "" {
		cfg.Server.Listen = ":8080"
	}
	if cfg.Auth.CookieName == "" {
		cfg.Auth.CookieName = "auth_token"
	}
	if cfg.Auth.CookieMaxAge == 0 {
		cfg.Auth.CookieMaxAge = 24 * time.Hour
	}
	if cfg.Auth.TokenDuration == 0 {
		cfg.Auth.TokenDuration = 24 * time.Hour
	}

	return &cfg, nil
}
