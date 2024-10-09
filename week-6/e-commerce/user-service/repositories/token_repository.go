package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/partridge1307/service-ctx/core"
	"github.com/redis/go-redis/v9"
)

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID int, deviceID, token string, expiration int) error
	GetRefreshToken(ctx context.Context, userID int, deviceID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID int, deviceID string) error
	DeleteAllRefreshToken(ctx context.Context, userID int) error
}

type redisRepo struct {
	db *redis.Client
}

func NewTokenRepo(db *redis.Client) TokenRepository {
	return &redisRepo{
		db,
	}
}

func (repo *redisRepo) SetRefreshToken(ctx context.Context, userID int, deviceID, token string, expiration int) error {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)

	err := repo.db.Set(ctx, key, token, time.Second*time.Duration(expiration)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *redisRepo) GetRefreshToken(ctx context.Context, userID int, deviceID string) (string, error) {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)

	token, err := repo.db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (repo *redisRepo) DeleteRefreshToken(ctx context.Context, userID int, deviceID string) error {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)

	affected, err := repo.db.Del(ctx, key).Result()
	if affected == 0 {
		return core.ErrRecordNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (repo *redisRepo) DeleteAllRefreshToken(ctx context.Context, userID int) error {
	key := fmt.Sprintf("refresh_token:%d:*", userID)

	affected, err := repo.db.Del(ctx, key).Result()
	if affected == 0 {
		return core.ErrRecordNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
