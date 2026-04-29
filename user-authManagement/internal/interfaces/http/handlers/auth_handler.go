package handlers

import (
	"log/slog"
	"net/http"

	"github.com/lindritprekaj/user-authmanagement/internal/application/auth"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/dto"
)

type AuthHandler struct {
	logger          *slog.Logger
	register        *auth.RegisterUseCase
	login           *auth.LoginUseCase
	refresh         *auth.RefreshUseCase
	logout          *auth.LogoutUseCase
	requestReset    *auth.RequestPasswordResetUseCase
	confirmReset    *auth.ConfirmPasswordResetUseCase
}

func NewAuthHandler(
	logger *slog.Logger,
	register *auth.RegisterUseCase,
	login *auth.LoginUseCase,
	refresh *auth.RefreshUseCase,
	logout *auth.LogoutUseCase,
	requestReset *auth.RequestPasswordResetUseCase,
	confirmReset *auth.ConfirmPasswordResetUseCase,
) *AuthHandler {
	return &AuthHandler{
		logger: logger, register: register, login: login,
		refresh: refresh, logout: logout,
		requestReset: requestReset, confirmReset: confirmReset,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	out, err := h.register.Execute(r.Context(), auth.RegisterInput{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusCreated, dto.RegisterResponse{UserID: out.UserID, Email: out.Email})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	out, err := h.login.Execute(r.Context(), auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		IP:        r.Header.Get("X-Forwarded-For"),
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, dto.TokenResponse{
		AccessToken: out.AccessToken, AccessExpiresAt: out.AccessExpiresAt,
		RefreshToken: out.RefreshToken, RefreshExpiresAt: out.RefreshExpiresAt,
		TokenType: "Bearer",
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	out, err := h.refresh.Execute(r.Context(), auth.RefreshInput{RefreshToken: req.RefreshToken})
	if err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, dto.TokenResponse{
		AccessToken: out.AccessToken, AccessExpiresAt: out.AccessExpiresAt,
		RefreshToken: out.RefreshToken, RefreshExpiresAt: out.RefreshExpiresAt,
		TokenType: "Bearer",
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	if err := h.logout.Execute(r.Context(), auth.LogoutInput{RefreshToken: req.RefreshToken}); err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestPasswordResetRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	// Always 202 to avoid leaking which emails are registered.
	if err := h.requestReset.Execute(r.Context(), auth.RequestPasswordResetInput{Email: req.Email}); err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmPasswordResetRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		return
	}
	if err := h.confirmReset.Execute(r.Context(), auth.ConfirmPasswordResetInput{
		ResetToken:  req.ResetToken,
		NewPassword: req.NewPassword,
	}); err != nil {
		writeDomainError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
