package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

// RequestIdMiddleware Request ID 미들웨어
func RequestIdMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 존재 여부 확인
		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			// UUID 생성
			requestId = uuid.New().String()
			r.Header.Set("X-Request-ID", requestId)
		}
		w.Header().Set("X-Request-ID", requestId)
		next.ServeHTTP(w, r)
	}
}
