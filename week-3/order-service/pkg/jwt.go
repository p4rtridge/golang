package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	IssueAccessToken(ctx context.Context, id, sub string) (string, int, error)
	IssueRefreshToken(ctx context.Context, id, sub string) (string, int, error)
	ParseToken(ctx context.Context, tokenStr string) (*jwt.RegisteredClaims, error)
}

type jwtx struct {
	secret_key    []byte
	atExpireInSec int
	rtExpireInSec int
}

func NewJWT(secret_key string, atExpireInSec, rtExpireInSec int) JWT {
	return &jwtx{
		secret_key:    []byte(secret_key),
		atExpireInSec: atExpireInSec,
		rtExpireInSec: rtExpireInSec,
	}
}

func (r *jwtx) IssueAccessToken(ctx context.Context, id, sub string) (string, int, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(r.atExpireInSec))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(r.secret_key)
	if err != nil {
		return "", 0, err
	}

	return signedToken, r.atExpireInSec, nil
}

func (r *jwtx) IssueRefreshToken(ctx context.Context, id, sub string) (string, int, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(r.rtExpireInSec))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(r.secret_key)
	if err != nil {
		return "", 0, err
	}

	return signedToken, r.rtExpireInSec, nil
}

func (r *jwtx) ParseToken(ctx context.Context, tokenStr string) (*jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return r.secret_key, nil
	})
	if !token.Valid {
		return nil, err
	}

	return &claims, nil
}
