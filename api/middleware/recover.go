package middleware

import (
	"fmt"
	"net/http"
)

// RecoverMiddleware recover 미들웨어
func RecoverMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic:", r)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
