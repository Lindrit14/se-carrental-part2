// Package user defines the User aggregate: entity, repository contract,
// and domain errors. This package has no framework dependencies.
package user

import (
	"strings"
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User is the aggregate root.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Roles        []Role
	Verified     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// New creates a User with normalized email and default role.
// It does not hash the password — the caller passes an already-hashed value.
func New(id, email, passwordHash string, now time.Time) *User {
	return &User{
		ID:           id,
		Email:        NormalizeEmail(email),
		PasswordHash: passwordHash,
		Roles:        []Role{RoleUser},
		Verified:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NormalizeEmail returns a canonical form of the address.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// HasRole reports whether the user holds the given role.
func (u *User) HasRole(r Role) bool {
	for _, rr := range u.Roles {
		if rr == r {
			return true
		}
	}
	return false
}

// MarkVerified flips the verified flag.
func (u *User) MarkVerified(now time.Time) {
	u.Verified = true
	u.UpdatedAt = now
}

// ChangePassword replaces the password hash.
func (u *User) ChangePassword(hash string, now time.Time) {
	u.PasswordHash = hash
	u.UpdatedAt = now
}
