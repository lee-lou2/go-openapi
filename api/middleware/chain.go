package middleware

import "net/http"

type Middleware func(handler http.Handler) http.HandlerFunc

// Chain 미들웨어 체인 생성
func Chain(middlewares ...Middleware) Middleware {
	return func(handler http.Handler) http.HandlerFunc {
		for i := range middlewares {
			handler = middlewares[len(middlewares)-1-i](handler)
		}
		return handler.ServeHTTP
	}
}
