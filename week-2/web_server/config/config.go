package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP `yaml:"http"`
		PG   `yaml:"pg"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PG struct {
		URL string `env-required:"true" env:"PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadConfig(fmt.Sprintf("%s/config/config.yml", dir), cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
