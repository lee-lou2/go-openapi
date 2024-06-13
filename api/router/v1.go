package router

import (
	"github.com/gofiber/fiber/v3"
	authHandler "go-openapi/api/handler/auth"
	clientHandler "go-openapi/api/handler/client"
	userHandler "go-openapi/api/handler/user"
	"go-openapi/api/middleware"
	clientModel "go-openapi/model/client"
	"net/http"
)

func V1Router(app fiber.Router) {
	v1 := app.Group("/v1")
	{
		user := v1.Group("/user")
		{
			// 이메일 검증
			user.Patch("/verify/:code", userHandler.VerifyCodeHandler)
			// 비밀번호 재설정 코드 전송
			user.Post("/password", userHandler.SendPasswordResetCodeHandler)
			// 비밀번호 재설정
			user.Patch("/password", userHandler.ResetPasswordHandler)
		}
		auth := v1.Group("/auth")
		{
			// Client 생성
			auth.Post("/client", authHandler.CreateClientHandler, middleware.AuthMiddleware, middleware.PermissionMiddleware(clientModel.ScopeWriteClient))
			// Client 조회
			auth.Get("/client", authHandler.GetClientsHandler, middleware.AuthMiddleware, middleware.PermissionMiddleware(clientModel.ScopeReadClient))
			// Client 삭제
			auth.Delete("/client/:id", authHandler.DeleteClientHandler, middleware.AuthMiddleware, middleware.PermissionMiddleware(clientModel.ScopeWriteClient))
			// 로그인
			auth.Post("/login", authHandler.LoginHandler)
			// 로그아웃
			auth.Post("/logout", authHandler.LogoutHandler, middleware.AuthMiddleware)

			// 토큰 발급
			auth.Post("/token", authHandler.CreateTokenHandler)
			// 토큰 갱신
			auth.Post("/token/refresh", authHandler.RefreshTokenHandler)
		}
		client := v1.Group("/client")
		{
			// 내 정보 조회
			client.Get("/me", clientHandler.GetMeHandler, middleware.AuthMiddleware, middleware.LimitPerSecondMiddleware(clientModel.ScopeReadMe))
		}
	}
}

// V1Router2 라우터 설정
func V1Router2() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user", userHandler.CreateUserHandler)
	mux.HandleFunc("POST /user/verify", userHandler.SendVerifyCodeHandler)
	return mux
}
