package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
)

// Authn requires a valid Bearer token. On success it injects user_id + roles
// into the request context.
func Authn(tokenSvc ports.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				writeUnauthorized(w)
				return
			}
			tok := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
			claims, err := tokenSvc.ParseAccessToken(tok)
			if err != nil {
				writeUnauthorized(w)
				return
			}
			ctx := WithUser(r.Context(), claims.UserID, claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    "unauthorized",
		"message": "missing or invalid access token",
	})
}
