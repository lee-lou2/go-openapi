package middleware

import (
	"net/http"
)

func PermissionMiddleware(scope string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scopes := r.Context().Value("scopes").([]string)
		for _, s := range scopes {
			if s == scope {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}
