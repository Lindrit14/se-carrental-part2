// Package handlers implements the HTTP endpoints. Handlers translate
// between dto and use-case input/output and never contain business logic.
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/dto"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/middleware"
	"github.com/lindritprekaj/user-authmanagement/pkg/apperror"
)

// validate is module-private; share one instance.
var validate = validator.New(validator.WithRequiredStructEnabled())

// decodeAndValidate reads JSON into v and validates struct tags.
// On failure it writes a 400 and returns a non-nil error so the caller can return early.
func decodeAndValidate(w http.ResponseWriter, r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, dto.ErrorResponse{
			Code: "invalid_body", Message: "request body is not valid JSON",
		})
		return err
	}
	if err := validate.Struct(v); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			details := make(map[string]string, len(verrs))
			for _, fe := range verrs {
				details[fe.Field()] = fe.Tag()
			}
			writeError(w, http.StatusBadRequest, dto.ErrorResponse{
				Code: "validation_failed", Message: "request validation failed", Details: details,
			})
			return err
		}
		writeError(w, http.StatusBadRequest, dto.ErrorResponse{Code: "invalid_request", Message: err.Error()})
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, body dto.ErrorResponse) {
	writeJSON(w, status, body)
}

// writeDomainError maps a domain error to its HTTP representation.
func writeDomainError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	status, env := apperror.Map(err)
	if status >= 500 {
		logger.LogAttrs(r.Context(), slog.LevelError, "handler_error",
			slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			slog.Any("error", err),
		)
	}
	writeJSON(w, status, dto.ErrorResponse{Code: env.Code, Message: env.Message})
}
