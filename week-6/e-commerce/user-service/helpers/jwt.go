package helpers

import (
	"context"
	"fmt"
	"time"

	jwtx "github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	IssueAccessToken(ctx context.Context, sub, id string) (string, int, error)
	IssueRefreshToken(ctx context.Context, sub, id string) (string, int, error)
	ParseToken(ctx context.Context, tokenStr string) (*ParsedToken, error)
}

type ParsedToken struct {
	Issuer    string     `json:"issuer,omitempty"`
	Subject   string     `json:"subject,omitempty"`
	Audience  []string   `json:"audience,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	NotBefore *time.Time `json:"not_before,omitempty"`
	IssuedAt  *time.Time `json:"issued_at,omitempty"`
	TokenID   string     `json:"token_id,omitempty"`
}

func (t ParsedToken) GetIssuer() string {
	return t.Issuer
}

func (t ParsedToken) GetSubject() string {
	return t.Subject
}

func (t ParsedToken) GetAudience() []string {
	return t.Audience
}

func (t ParsedToken) GetExpiresAt() *time.Time {
	return t.ExpiresAt
}

func (t ParsedToken) GetNotBefore() *time.Time {
	return t.NotBefore
}

func (t ParsedToken) GetIssuedAt() *time.Time {
	return t.IssuedAt
}

func (t ParsedToken) GetTokenID() string {
	return t.TokenID
}

type jwt struct {
	secret_key    []byte
	atExpireInSec int
	rtExpireInSec int
	audience      []string
}

func NewJWT(secretStr string, audience []string, atExpireInSec, rtExpireInSec int) JWT {
	return &jwt{
		secret_key:    []byte(secretStr),
		atExpireInSec: atExpireInSec,
		rtExpireInSec: rtExpireInSec,
		audience:      audience,
	}
}

func (t *jwt) IssueAccessToken(ctx context.Context, sub, id string) (string, int, error) {
	now := time.Now()

	claims := jwtx.RegisteredClaims{
		Subject:   sub,
		ID:        id,
		ExpiresAt: jwtx.NewNumericDate(now.Add(time.Second * time.Duration(t.atExpireInSec))),
		NotBefore: jwtx.NewNumericDate(now),
		IssuedAt:  jwtx.NewNumericDate(now),
		Audience:  jwtx.ClaimStrings{"e-commerce"},
	}

	signedToken, err := jwtx.NewWithClaims(jwtx.SigningMethodHS256, claims).SignedString(t.secret_key)
	if err != nil {
		return "", 0, err
	}

	return signedToken, t.atExpireInSec, nil
}

func (t *jwt) IssueRefreshToken(ctx context.Context, sub, id string) (string, int, error) {
	now := time.Now()

	claims := jwtx.RegisteredClaims{
		Subject:   sub,
		ID:        id,
		ExpiresAt: jwtx.NewNumericDate(now.Add(time.Second * time.Duration(t.rtExpireInSec))),
		NotBefore: jwtx.NewNumericDate(now),
		IssuedAt:  jwtx.NewNumericDate(now),
		Audience:  jwtx.ClaimStrings{"e-commerce"},
	}

	signedToken, err := jwtx.NewWithClaims(jwtx.SigningMethodHS256, claims).SignedString(t.secret_key)
	if err != nil {
		return "", 0, err
	}

	return signedToken, t.rtExpireInSec, nil
}

func (t *jwt) ParseToken(ctx context.Context, tokenStr string) (*ParsedToken, error) {
	var claims jwtx.RegisteredClaims

	token, err := jwtx.ParseWithClaims(tokenStr, &claims, func(token *jwtx.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtx.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secret_key, nil
	})

	if !token.Valid {
		return nil, err
	}

	var parsedToken ParsedToken

	// shallow copy
	parsedToken.Issuer = claims.Issuer
	parsedToken.Subject = claims.Subject
	parsedToken.Audience = claims.Audience
	parsedToken.ExpiresAt = &claims.ExpiresAt.Time
	parsedToken.NotBefore = &claims.NotBefore.Time
	parsedToken.IssuedAt = &claims.IssuedAt.Time
	parsedToken.TokenID = claims.ID

	return &parsedToken, nil
}
