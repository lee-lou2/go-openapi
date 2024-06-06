package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-openapi/config"
	"strings"
	"time"
)

// Token 토큰 구조체
type Token struct {
	JTI       string `json:"jti"`
	TokenType string `json:"token_type"`
	Sub       uint   `json:"sub"`
	SubType   string `json:"sub_type"`
	Scope     string `json:"scope"`
	Iat       int64  `json:"iat"`
	Exp       int64  `json:"exp"`
}

// GetTokenClaims 토큰 클레임
func GetTokenClaims(tokenString string) (*Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	var t Token
	t.JTI = claims["jti"].(string)
	t.TokenType = claims["token_type"].(string)
	t.Sub = uint(claims["sub"].(float64))
	t.SubType = claims["sub_type"].(string)
	t.Scope = claims["scope"].(string)
	t.Iat = int64(claims["iat"].(float64))
	t.Exp = int64(claims["exp"].(float64))
	return &t, nil
}

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
