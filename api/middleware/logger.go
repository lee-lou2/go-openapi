package middleware

import (
	"log"
	"net/http"
)

// LoggerMiddleware 로거 미들웨어
func LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
