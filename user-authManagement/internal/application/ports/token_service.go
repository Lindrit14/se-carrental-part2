package ports

import "time"

// AccessClaims is the application-facing view of a parsed JWT.
type AccessClaims struct {
	UserID    string
	Roles     []string
	ExpiresAt time.Time
}

// TokenService issues and parses access JWTs and produces the
// opaque random material used as a refresh token.
type TokenService interface {
	IssueAccessToken(userID string, roles []string) (token string, expiresAt time.Time, err error)
	ParseAccessToken(token string) (*AccessClaims, error)

	// NewRefreshToken returns a fresh, opaque, URL-safe random string.
	// The plaintext is given to the client; the hash is persisted.
	NewRefreshToken() (plaintext string, err error)
	HashRefreshToken(plaintext string) string
}
