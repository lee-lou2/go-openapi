package router

import (
	authHandler "go-openapi/api/handler/auth"
	clientHandler "go-openapi/api/handler/client"
	userHandler "go-openapi/api/handler/user"
	"go-openapi/api/middleware"
	clientModel "go-openapi/model/client"
	"net/http"
)

// V1Router 라우터 설정
func V1Router() *http.ServeMux {
	mux := http.NewServeMux()
	userMux := http.NewServeMux()
	{
		// 사용자 생성
		userMux.HandleFunc("POST /", userHandler.CreateUserHandler)

		// 사용자 확인 및 인증 코드 전송
		userMux.HandleFunc("POST /verify/", userHandler.SendVerifyCodeHandler)
		// 인증 코드 확인 및 사용자 업데이트
		userMux.HandleFunc("PATCH /verify/{code}/", userHandler.VerifyCodeHandler)

		// 비밀번호 재설정 코드 전송
		userMux.HandleFunc("POST /password/", userHandler.SendPasswordResetCodeHandler)
		// 비밀번호 재설정
		userMux.HandleFunc("PATCH /password/", userHandler.ResetPasswordHandler)
	}
	authMux := http.NewServeMux()
	{
		// 토큰 발급
		authMux.HandleFunc("POST /token/", authHandler.CreateTokenHandler)
		// 토큰 갱신
		authMux.HandleFunc("POST /token/refresh/", authHandler.RefreshTokenHandler)

		// Client 생성
		authMux.Handle("POST /client/", middleware.AuthMiddleware(middleware.PermissionMiddleware(clientModel.ScopeWriteClient, http.HandlerFunc(authHandler.CreateClientHandler))))
		// Client 조회
		authMux.Handle("GET /client/", middleware.AuthMiddleware(middleware.PermissionMiddleware(clientModel.ScopeReadClient, http.HandlerFunc(authHandler.GetClientsHandler))))
		// // Client 삭제
		authMux.Handle("DELETE /client/{id}/", middleware.AuthMiddleware(middleware.PermissionMiddleware(clientModel.ScopeWriteClient, http.HandlerFunc(authHandler.DeleteClientHandler))))

		// 로그인
		authMux.HandleFunc("POST /login/", authHandler.LoginHandler)
		// 로그아웃
		authMux.HandleFunc("POST /logout/", middleware.AuthMiddleware(http.HandlerFunc(authHandler.LogoutHandler)))
	}
	clientMux := http.NewServeMux()
	{
		// 내 정보 조회
		clientMux.HandleFunc("GET /me/", middleware.AuthMiddleware(middleware.LimitPerSecondMiddleware(clientModel.ScopeReadMe, http.HandlerFunc(clientHandler.GetMeHandler))))
	}
	mux.Handle("/user/", http.StripPrefix("/user", userMux))
	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))
	return mux
}
