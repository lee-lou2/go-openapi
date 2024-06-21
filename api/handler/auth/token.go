package auth

import (
	"fmt"
	"go-openapi/api/parser"
	"go-openapi/api/render"
	"go-openapi/api/validation"
	authInternal "go-openapi/internal/auth"
	clientModel "go-openapi/model/client"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// CreateTokenHandler 클라이언트 자격 증명 핸들러
func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 요청 데이터 유효성 검사
	var body struct {
		Scope     string `json:"scope"`
		GrantType string `json:"grant_type"`
	}
	if err := parser.JSON(r, &body); err != nil {
		render.JSON(w, fiber.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !validation.ValidateClientCredentials(body.Scope, body.GrantType) {
		render.JSON(w, fiber.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}
	token, err := authInternal.GetTokenFromClient(r, body.Scope)
	if err != nil {
		render.JSON(w, fiber.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, fiber.StatusOK, map[string]string{
		"tokenType":   "Bearer",
		"accessToken": token,
		"expiresIn":   strconv.Itoa(60 * 15),
		"scope":       body.Scope,
	})
}

// RefreshTokenHandler 토큰 갱신 핸들러
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := parser.JSON(r, &body); err != nil {
		render.JSON(w, fiber.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !validation.ValidateRefreshToken(body.GrantType, body.RefreshToken) {
		render.JSON(w, fiber.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}
	token, err := authInternal.GetTokenFromRefreshToken(body.RefreshToken)
	if err != nil {
		render.JSON(w, fiber.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, fiber.StatusOK, map[string]string{
		"tokenType":   "Bearer",
		"accessToken": token,
		"expiresIn":   strconv.Itoa(60 * 15),
		"scope":       fmt.Sprintf("%s %s", clientModel.ScopeWriteClient, clientModel.ScopeReadClient),
	})
}
