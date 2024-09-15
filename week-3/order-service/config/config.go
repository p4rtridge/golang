package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type JWTCfg struct {
	SecretKey     string `env-required:"true" env:"JWT_SECRET_KEY"`
	ATExpireInSec int    `env:"JWT_ACCESS_TOKEN_EXPIRE_IN_SEC" env-default:"604800"`   // 60 * 60 * 24 * 7
	RTExpireInSec int    `env:"JWT_REFRESH_TOKEN_EXPIRE_IN_SEC" env-default:"2592000"` // 60 * 60 * 24 * 30
}

type PGCfg struct {
	URL string `env-required:"true" env:"PG_URL"`
}

type RDCfg struct {
	URL string `env-required:"true" env:"REDIS_URL"`
}

type Config struct {
	PGCfg
	RDCfg
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
		pool, err := pgxpool.New(ctx, cfg.PGCfg.URL)

		if err == nil && pool != nil {
			return pool
		}

		e = err

		time.Sleep(time.Second * 2)
	}

	log.Fatalln(fmt.Errorf("postgres connect error: %v", e))

	return nil
}

func ConnectToRedis(cfg *Config) *redis.Client {
	opts, err := redis.ParseURL(cfg.RDCfg.URL)
	if err != nil {
		log.Fatalln(fmt.Errorf("redis parse config error: %v", err))
	}

	return redis.NewClient(opts)
}
