package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v3"
	authCmd "go-openapi/cmd/auth"
	"go-openapi/config"
	clientModel "go-openapi/model/client"
	"strings"
)

// getClientKeys 클라이언트 키 가져오기
func getClientKeys(c fiber.Ctx) (string, string, error) {
	authHeader := c.Get("Authorization")
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

// CreateTokenHandler 클라이언트 자격 증명 핸들러
func CreateTokenHandler(c fiber.Ctx) error {
	// 요청 데이터 유효성 검사
	body := struct {
		Scope     string `json:"scope"`
		GrantType string `json:"grant_type"`
	}{}
	if err := c.Bind().JSON(&body); err != nil {
		return err
	}
	if body.GrantType == "" || body.Scope == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "scope is required",
		})
	}
	if body.GrantType != "client_credentials" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unsupported_grant_type",
		})
	}
	// 클라이언트 유효성 검사
	clientId, clientSecret, err := getClientKeys(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// 데이터 조회
	db := config.GetDB()
	var client clientModel.Client
	result := db.Where("client_id = ? AND client_secret = ?", clientId, clientSecret).First(&client)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// 권한 검증
	scopes := strings.Split(client.Scope, " ")
	bodyScopes := strings.Split(body.Scope, " ")
	isExist := false
	for _, bodyScope := range bodyScopes {
		for _, clientScope := range scopes {
			if bodyScope == clientScope {
				isExist = true
				break
			}
		}
		if !isExist {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid_scope",
			})
		}
	}
	// 토큰 생성(15분)
	exp := 60 * 15
	token, err := authCmd.CreateToken("access", client.ID, exp, scopes...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"tokenType":   "Bearer",
		"accessToken": token,
		"expiresIn":   exp,
		"scope":       body.Scope,
	})
}
