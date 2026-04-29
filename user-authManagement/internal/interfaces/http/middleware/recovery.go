package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery converts panics into 500 responses and logs the stack.
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.LogAttrs(r.Context(), slog.LevelError, "http_panic",
						slog.String("request_id", RequestIDFromContext(r.Context())),
						slog.Any("panic", rec),
						slog.String("stack", string(debug.Stack())),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"code":"internal_error","message":"internal server error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
