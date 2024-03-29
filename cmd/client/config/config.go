// Package config содержит конфигурацию клиента.
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	flag "github.com/spf13/pflag"
)

type (
	// Config основная конфигурация.
	Config struct {
		App     `yaml:"app"`
		Storage `yaml:"storage"`
		GRPC    `yaml:"grpc"`
	}

	// App информация о приложении.
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// Storage настройки файлового хранилища.
	Storage struct {
		Path string `env-required:"true" yaml:"path" env:"STORAGE_PATH"`
	}

	// GRPC настройки gRPC.
	GRPC struct {
		Address string `env-required:"true" yaml:"address" env:"GRPC_ADDRESS"`
		Port    string `env-required:"true" yaml:"port"    env:"GRPC_PORT"`
	}
)

// New создаёт объект Config.
func New() (*Config, error) {
	cfg := &Config{}

	var cfgFile string
	flag.StringVarP(&cfgFile, "config", "c", "", "path to config file (*.yaml)")
	flag.Parse()

	err := cleanenv.ReadConfig(cfgFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
