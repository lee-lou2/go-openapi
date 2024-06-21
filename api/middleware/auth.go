package middleware

import (
	"context"
	authPkg "go-openapi/pkg/auth"
	"net/http"
	"strings"
)

// AuthMiddleware 사용자 인증 미들웨어
func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		// 토큰에서 데이터 조회
		claims, err := authPkg.GetTokenClaims(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// 사용자 정보를 context에 저장
		ctx := r.Context()
		ctx = context.WithValue(ctx, claims.SubType, claims.Sub)
		scopes := strings.Split(claims.Scope, " ")
		ctx = context.WithValue(ctx, "scopes", scopes)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
