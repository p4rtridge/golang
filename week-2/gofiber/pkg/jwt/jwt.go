package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partridge1307/gofiber/entity"
)

func Sign(user *entity.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	return claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func Verify(tokenStr string) (*entity.User, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	user := &entity.User{}

	if userID, ok := claims["id"].(float64); ok {
		user.SetID(int(userID))
	} else {
		return nil, errors.New("poison claims")
	}

	if username, ok := claims["sub"].(string); ok {
		user.SetUsername(username)
	} else {
		return nil, errors.New("poison claims")
	}

	return user, nil
}
