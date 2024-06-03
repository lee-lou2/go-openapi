package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go-openapi/config"
	"time"
)

// CreateTokenSet 토큰 셋 생성
func CreateTokenSet(userId uint) (accessToken, refreshToken string, err error) {
	accessToken, err = CreateToken("access", userId, 3600)
	if err != nil {
		return
	}
	refreshToken, err = CreateToken("refresh", userId, 86400)
	if err != nil {
		return
	}
	return accessToken, refreshToken, nil
}

// CreateToken 토큰 생성
func CreateToken(tokenType string, userId uint, exp int) (string, error) {
	claims := jwt.MapClaims{
		"token_type": tokenType,
		"user":       userId,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(time.Second * time.Duration(exp)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}
