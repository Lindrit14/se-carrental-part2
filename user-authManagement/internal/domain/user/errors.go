package user

import "errors"

var (
	ErrNotFound          = errors.New("user not found")
	ErrEmailTaken        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotVerified       = errors.New("email not verified")
	ErrWeakPassword      = errors.New("password does not meet policy")
	ErrInvalidEmail      = errors.New("invalid email address")
)
