package router

import "net/http"

// BaseRouter 기본 라우터
func BaseRouter() *http.ServeMux {
	// Health check
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", http.StripPrefix("/", fs))
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	return mux
}
