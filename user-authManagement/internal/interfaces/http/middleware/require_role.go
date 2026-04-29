package middleware

import (
	"encoding/json"
	"net/http"
)

// RequireRole returns 403 if the request context does not carry the given role.
// Must run after Authn, which populates the context with the JWT roles.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, r2 := range RolesFromContext(r.Context()) {
				if r2 == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"code":    "forbidden",
				"message": "insufficient role",
			})
		})
	}
}
