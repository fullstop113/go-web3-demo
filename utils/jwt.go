package utils

import (
	"errors"
	"time"

	"github.com/fullstop113/go-web3-demo/config"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func ParseToken(tokenString string) (*Claims, error) {
	secret := config.LoadJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	Claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return Claims, nil
}

func GenerateToken(userID uint, username string) (string, time.Time, error) {
	secret := config.LoadJWTSecret()
	expireAT := time.Now().Add(24 * time.Hour)

	claims := Claims{
		UserID: userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAT),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", expireAT, err
	}
	return tokenString, expireAT, nil
}