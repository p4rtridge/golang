package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Name   string `env:"APP_NAME" env-default:"Gofiber"`
	Header string `env:"APP_HEADER" env-default:"Fiber"`
}

type PG struct {
	URL         string `env-required:"true" env:"PG_URL"`
	MAX_RETRIES int    `env:"PG_MAX_RETRIES" env-default:"5"`
}

type Sign struct {
	SECRET_KEY string `env-required:"true" env:"SECRET_KEY"`
}

type Config struct {
	App
	PG
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
