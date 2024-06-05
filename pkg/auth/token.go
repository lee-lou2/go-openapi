package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-openapi/config"
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
