package router

import (
	"github.com/gofiber/fiber/v3"
	authHandler "go-openapi/api/handler/auth"
	userHandler "go-openapi/api/handler/user"
	"go-openapi/api/middleware"
)

func V1Router(app fiber.Router) {
	v1 := app.Group("/v1")
	{
		user := v1.Group("/user")
		{
			// 사용자 생성
			user.Post("", userHandler.CreateUserHandler)
			// 인증번호 재전송
			user.Post("/verify", userHandler.SendVerifyCodeHandler)
			// 이메일 검증
			user.Patch("/verify/:code", userHandler.VerifyCodeHandler)
			// 비밀번호 재설정 코드 전송
			user.Post("/password", userHandler.SendPasswordResetCodeHandler)
			// 비밀번호 재설정
			user.Patch("/password/:code", userHandler.ResetPasswordHandler)
		}
		auth := v1.Group("/auth")
		{
			// Client 생성
			auth.Post("/client", authHandler.CreateClientHandler, middleware.AuthMiddleware, middleware.PermissionMiddleware("write:client"))
			// Client 조회
			auth.Get("/client", authHandler.GetClientsHandler, middleware.AuthMiddleware, middleware.PermissionMiddleware("read:client"))
			// Client 삭제
			auth.Delete("/client/:id", authHandler.DeleteClientHandler, middleware.AuthMiddleware, middleware.PermissionMiddleware("write:client"))
			// 로그인
			auth.Post("/login", authHandler.LoginHandler)
			// 로그아웃
			auth.Post("/logout", authHandler.LogoutHandler, middleware.AuthMiddleware)

			// 토큰 발급
			auth.Post("/token", authHandler.CreateTokenHandler)
		}
	}
}
