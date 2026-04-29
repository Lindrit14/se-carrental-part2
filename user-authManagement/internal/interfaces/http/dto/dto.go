// Package dto contains the HTTP request/response shapes. Domain entities
// never cross the HTTP boundary — handlers translate them via these types.
package dto

import "time"

// --- Requests --------------------------------------------------------

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmPasswordResetRequest struct {
	ResetToken  string `json:"reset_token"  validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=12"`
}

type UpdateProfileRequest struct {
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

// --- Responses -------------------------------------------------------

type RegisterResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenResponse struct {
	AccessToken      string    `json:"access_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	TokenType        string    `json:"token_type"`
}

type ProfileResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Verified  bool     `json:"verified"`
	CreatedAt string   `json:"created_at"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// --- Admin -----------------------------------------------------------

type AdminUserResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Verified  bool     `json:"verified"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type AdminUserListResponse struct {
	Items  []AdminUserResponse `json:"items"`
	Total  int64               `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

type AdminSetRolesRequest struct {
	Roles []string `json:"roles" validate:"required,min=1,dive,oneof=user admin"`
}
