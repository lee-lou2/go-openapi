package api

import (
	"go-openapi/api/middleware"
	"go-openapi/api/router"
	"net/http"
)

// Server 서버 설정
func Server() http.HandlerFunc {
	mux := http.NewServeMux()
	base := router.BaseRouter()
	mux.Handle("/", base)
	v1 := router.V1Router()
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	chain := middleware.Chain(
		middleware.LoggerMiddleware,
		middleware.RecoverMiddleware,
		middleware.CORSMiddleware,
		middleware.RequestIdMiddleware,
	)
	return chain(mux)
}
