// Package apperror centralises the mapping from domain errors to
// HTTP status codes and a stable error envelope.
package apperror

import (
	"errors"
	"net/http"

	"github.com/lindritprekaj/user-authmanagement/internal/application/admin"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/token"
	"github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

type Envelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Map translates an error to (status, envelope). Unknown errors map to 500
// with a generic message — internal details stay in the log.
func Map(err error) (int, Envelope) {
	switch {
	case errors.Is(err, user.ErrNotFound):
		return http.StatusNotFound, Envelope{Code: "user_not_found", Message: "user not found"}
	case errors.Is(err, user.ErrEmailTaken):
		return http.StatusConflict, Envelope{Code: "email_taken", Message: "email already registered"}
	case errors.Is(err, user.ErrInvalidCredentials):
		return http.StatusUnauthorized, Envelope{Code: "invalid_credentials", Message: "invalid credentials"}
	case errors.Is(err, user.ErrNotVerified):
		return http.StatusForbidden, Envelope{Code: "email_not_verified", Message: "email not verified"}
	case errors.Is(err, user.ErrWeakPassword):
		return http.StatusBadRequest, Envelope{Code: "weak_password", Message: "password does not meet policy"}
	case errors.Is(err, user.ErrInvalidEmail):
		return http.StatusBadRequest, Envelope{Code: "invalid_email", Message: "invalid email address"}
	case errors.Is(err, token.ErrNotFound), errors.Is(err, token.ErrRevoked), errors.Is(err, token.ErrExpired):
		return http.StatusUnauthorized, Envelope{Code: "invalid_token", Message: "invalid or expired token"}
	case errors.Is(err, admin.ErrLastAdmin):
		return http.StatusConflict, Envelope{Code: "last_admin", Message: "cannot remove the last admin"}
	default:
		return http.StatusInternalServerError, Envelope{Code: "internal_error", Message: "internal server error"}
	}
}
