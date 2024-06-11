package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/validation"
	authInternal "go-openapi/internal/auth"
	clientModel "go-openapi/model/client"
)

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
	if !validation.ValidateClientCredentials(body.Scope, body.GrantType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_request",
		})
	}
	token, err := authInternal.GetTokenFromClient(c, body.Scope)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"tokenType":   "Bearer",
		"accessToken": token,
		"expiresIn":   60 * 15,
		"scope":       body.Scope,
	})
}

// RefreshTokenHandler 토큰 갱신 핸들러
func RefreshTokenHandler(c fiber.Ctx) error {
	body := struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}{}
	if err := c.Bind().JSON(&body); err != nil {
		return err
	}
	if !validation.ValidateRefreshToken(body.GrantType, body.RefreshToken) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_request",
		})
	}
	token, err := authInternal.GetTokenFromRefreshToken(body.RefreshToken)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"tokenType":   "Bearer",
		"accessToken": token,
		"expiresIn":   60 * 15,
		"scope":       fmt.Sprintf("%s %s", clientModel.ScopeWriteClient, clientModel.ScopeReadClient),
	})
}
