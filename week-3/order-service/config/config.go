package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JWTCfg struct {
	SecretKey        string `env-required:"true" env:"JWT_SECRET_KEY"`
	ExpireTokenInSec int    `env:"JWT_EXPIRE_TOKEN_IN_SEC" env-default:"604800"` // 60 * 60 * 24 * 7
}

type PGCfg struct {
	URL string `env-required:"true" env:"PG_URL"`
}

type Config struct {
	PGCfg
	JWTCfg
}

func NewConfig() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalln(fmt.Errorf("read config error: %v", err))
	}

	return &cfg
}

func ConnectToPostgres(cfg *Config) *pgxpool.Pool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var e error

	for i := 0; i < 5; i++ {
		pool, err := pgxpool.New(ctx, cfg.URL)

		if err == nil && pool != nil {
			return pool
		}

		e = err

		time.Sleep(time.Second * 2)
	}

	log.Fatalln(fmt.Errorf("postgres connect error: %v", e))

	return nil
}
