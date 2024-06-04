package middleware

import (
	"github.com/gofiber/fiber/v3"
	authPkg "go-openapi/pkg/auth"
	"strings"
)

// AuthMiddleware 사용자 인증 미들웨어
func AuthMiddleware(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	// 토큰에서 데이터 조회
	claims, err := authPkg.GetTokenClaims(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	scopes := strings.Split(claims.Scope, " ")
	fiber.Locals[uint](c, "user", claims.User)
	fiber.Locals[[]string](c, "scopes", scopes)
	return c.Next()
}
