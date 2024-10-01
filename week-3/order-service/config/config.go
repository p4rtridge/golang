package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

type AWSCfg struct {
	EndPoint  string `env-required:"true" env:"AWS_S3_ENDPOINT"`
	Region    string `env-required:"true" env:"AWS_S3_REGION"`
	AccessKey string `env-required:"true" env:"AWS_S3_ACCESS_KEY"`
	SecretKey string `env-required:"true" env:"AWS_S3_SECRET_KEY"`
}

type Config struct {
	PGCfg
	RDCfg
	AWSCfg
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

	pgCfg, err := pgxpool.ParseConfig(cfg.PGCfg.URL)
	if err != nil {
		log.Fatalln(fmt.Errorf("postgres parse config error: %v", err))
	}

	// to handle 10000 concurrent users
	pgCfg.MaxConns = 100

	for i := 0; i < 5; i++ {
		pool, err := pgxpool.NewWithConfig(ctx, pgCfg)

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

	// to handle 10000 concurrent users
	opts.PoolSize = 100

	return redis.NewClient(opts)
}

func ConnectToAWS(cfg *Config) *s3.Client {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	aws_cfg, err := awsCfg.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalln(fmt.Errorf("aws parse config error: %v", err))
	}

	client := s3.NewFromConfig(aws_cfg, func(o *s3.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.AWSCfg.SecretKey, "")
		o.BaseEndpoint = aws.String(cfg.EndPoint)
		o.Region = cfg.Region
		o.UsePathStyle = true
	})

	return client
}
