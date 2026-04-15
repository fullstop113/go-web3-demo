package utils

import (
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

	Claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return Claims, nil
}