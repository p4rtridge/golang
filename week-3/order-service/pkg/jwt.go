package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtx struct {
	secret_key       []byte
	expireTokenInSec int
}

func NewJWT(secret_key string, expireTokenInSec int) *jwtx {
	return &jwtx{
		secret_key:       []byte(secret_key),
		expireTokenInSec: expireTokenInSec,
	}
}

func (r *jwtx) IssueToken(ctx context.Context, id, sub string) (string, int, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(r.expireTokenInSec))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(r.secret_key)
	if err != nil {
		return "", 0, err
	}

	return token, r.expireTokenInSec, nil
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
