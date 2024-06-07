package auth

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	authPkg "go-openapi/pkg/auth"
	"go-openapi/pkg/utils"
)

// GetTokenFromLogin 로그인으로 토큰 발급
func GetTokenFromLogin(email string, password string) (accessToken string, refreshToken string, err error) {
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	hashedEmail := utils.SHA256Email(email)
	if err := db.Where("hashed_email = ?", hashedEmail).Where("is_verified = true").First(&user).Error; err != nil {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "invalid user")
	}
	// 비밀번호 확인
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "invalid user")
	}
	// 토큰 생성(사용자 로그인의 경우 클라이언트 관리만 가능)
	accessToken, refreshToken, err = authPkg.CreateTokenSet(user.ID, "user", "read:client", "write:client")
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return accessToken, refreshToken, nil
}
