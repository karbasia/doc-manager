package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/karbasia/doc-manager/internal/env"
)

type Config struct {
	Env      string // development | staging | production
	LogLevel string

	HTTP HTTPConfig
	DB   DBConfig

	JWTSecret string
	JWTIssuer string
}

type HTTPConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Path         string
	MaxReadConns int
	BusyTimeout  time.Duration
}

func (c *Config) IsProduction() bool { return c.Env == "production" }

func (h HTTPConfig) Addr() string { return fmt.Sprintf(":%d", h.Port) }

func Load() (*Config, error) {

	jwtSecret, err := env.GetRequiredString("JWT_SECRET")
	if err != nil {
		return nil, errors.New("JWT Secret not set")
	}

	jwtIssuer, err := env.GetRequiredString("JWT_ISSUER")
	if err != nil {
		return nil, errors.New("JWT Issuer not set")
	}

	cfg := &Config{
		Env:      env.GetString("APP_ENV", "development"),
		LogLevel: env.GetString("LOG_LEVEL", "info"),
		HTTP: HTTPConfig{
			Port:            env.GetInt("PORT", 8080),
			ReadTimeout:     env.GetDuration("HTTP_READ_TIMEOUT", 10),
			WriteTimeout:    env.GetDuration("HTTP_WRITE_TIMEOUT", 10),
			ShutdownTimeout: env.GetDuration("HTTP_SHUTDOWN_TIMEOUT", 15),
		},
		DB: DBConfig{
			Path:         env.GetString("DB_PATH", "./data/app.db"),
			MaxReadConns: env.GetInt("DB_MAX_READ_CONNS", 4),
			BusyTimeout:  env.GetDuration("DB_BUSY_TIMEOUT", 5),
		},
		JWTSecret: jwtSecret,
		JWTIssuer: jwtIssuer,
	}

	return cfg, nil
}

func LoadAppConfig() *Config {
	cfg, err := Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return cfg
}
