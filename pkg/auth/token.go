package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-openapi/config"
)

// ValidateToken 토큰 유효성 검사
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetEnv("JWT_SECRET")), nil
	}); err != nil {
		return nil, err
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		} else {
			return nil, fmt.Errorf("invalid token")
		}
	}
}

// GetUser 사용자 조회
func GetUser(token string) (uint, error) {
	claims, err := ValidateToken(token)
	if err != nil {
		return 0, fmt.Errorf("unauthorized")
	}
	user, ok := claims["user"]
	if !ok {
		return 0, fmt.Errorf("user claim not found")
	}
	return uint(user.(float64)), nil
}
