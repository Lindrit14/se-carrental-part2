// Package token defines the RefreshToken aggregate. The plaintext token
// is never stored — only its hash. Rotation creates a new token and revokes
// the previous one in the same transaction.
package token

import "time"

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	IssuedAt  time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
	// ReplacedBy points to the rotated successor, if any.
	ReplacedBy string
}

// IsActive reports whether the token can still be used.
func (t *RefreshToken) IsActive(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

// Revoke marks the token as revoked at the given time.
func (t *RefreshToken) Revoke(now time.Time, replacedBy string) {
	t.RevokedAt = &now
	t.ReplacedBy = replacedBy
}
