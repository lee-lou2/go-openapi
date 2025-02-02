package auth

import (
	"encoding/base64"
	"fmt"
	"go-openapi/config"
	clientModel "go-openapi/model/client"
	authPkg "go-openapi/pkg/auth"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// getClient 클라이언트 키 가져오기
func getClient(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", fmt.Errorf("authorization header is required")
	}
	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", "", fmt.Errorf("invalid authorization header")
	}
	decoded, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		return "", "", fmt.Errorf("invalid authorization header")
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid authorization header")
	}
	return parts[0], parts[1], nil
}

// GetTokenFromClient 클라이언트 자격 증명으로 토큰 발급
func GetTokenFromClient(r *http.Request, scope string) (string, error) {
	// 클라이언트 유효성 검사
	clientId, clientSecret, err := getClient(r)
	if err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// 데이터 조회
	db := config.GetDB()
	var client clientModel.Client
	result := db.Where("client_id = ? AND client_secret = ?", clientId, clientSecret).First(&client)
	if result.RowsAffected == 0 {
		return "", fiber.NewError(fiber.StatusUnauthorized, "invalid client")
	}
	// 권한 검증
	scopes := strings.Split(client.Scope, " ")
	bodyScopes := strings.Split(scope, " ")
	isExist := false
	for _, bodyScope := range bodyScopes {
		for _, clientScope := range scopes {
			if bodyScope == clientScope {
				isExist = true
				break
			}
		}
		if !isExist {
			return "", fiber.NewError(fiber.StatusUnauthorized, "invalid scope")
		}
	}
	// 토큰 생성(15분)
	exp := 60 * 15
	token, err := authPkg.CreateToken("access", client.ID, "client", exp, scopes...)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return token, nil
}

// GetTokenFromRefreshToken 토큰 갱신
func GetTokenFromRefreshToken(refreshToken string) (string, error) {
	claims, err := authPkg.GetTokenClaims(refreshToken)
	if err != nil || claims.TokenType != "refresh" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	// 토큰 생성(15분)
	exp := 60 * 15
	token, err := authPkg.CreateToken("access", claims.Sub, "user", exp, clientModel.ScopeWriteClient, clientModel.ScopeReadClient)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return token, nil
}
