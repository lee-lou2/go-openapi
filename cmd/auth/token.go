package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-openapi/config"
	"strings"
	"time"
)

// CreateTokenSet 토큰 셋 생성
func CreateTokenSet(sub uint, subType string, scopes ...string) (accessToken, refreshToken string, err error) {
	accessToken, err = CreateToken("access", sub, subType, 3600, scopes...)
	if err != nil {
		return
	}
	refreshToken, err = CreateToken("refresh", sub, subType, 86400, scopes...)
	if err != nil {
		return
	}
	return accessToken, refreshToken, nil
}

// CreateToken 토큰 생성
func CreateToken(tokenType string, sub uint, subType string, exp int, scopes ...string) (string, error) {
	scope := strings.Join(scopes, " ")
	tokenId := uuid.New().String()
	claims := jwt.MapClaims{
		"jti":        tokenId,
		"token_type": tokenType,
		"sub":        sub,
		"sub_type":   subType,
		"scope":      scope,
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
