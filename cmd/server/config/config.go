// Package config содержит конфигурацию сервера.
package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config основная конфигурация.
	Config struct {
		App   `yaml:"app"`
		PG    `yaml:"postgres"`
		GRPC  `yaml:"grpc"`
		Token `yaml:"token"`
	}

	// App информация о приложении.
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// PG подключение к БД Postgres.
	PG struct {
		MaxOpen      int    `yaml:"pool_max"      env:"PG_POOL_MAX"`
		ConnAttempts int    `yaml:"conn_attempts" env:"PG_CONN_ATTEMPTS"`
		URL          string `env-required:"true"  env:"PG_URL"`
	}

	// GRPC настройки gRPC.
	GRPC struct {
		Port string `yaml:"port" env:"GRPC_PORT"`
	}

	// Token настройки для формирования токена.
	Token struct {
		Key            string        `env-required:"true"                        env:"TOKEN_KEY"`
		AccessDuration time.Duration `env-required:"true" yaml:"access_duration" env:"TOKEN_ACCESS_DURATION"`
	}
)

// New создаёт объект Config.
func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./cmd/server/config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
