package router

import "net/http"

// BaseRouter 기본 라우터
func BaseRouter() *http.ServeMux {
	// Health check
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./views"))
	mux.Handle("/", http.StripPrefix("/", fs))
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		// r.JSON(w, http.StatusOK, map[string]string{"message": "OK"})
	}))
	return mux
}
