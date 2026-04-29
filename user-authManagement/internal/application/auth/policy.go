package auth

import (
	"time"
	"unicode/utf8"

	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
)

const rfc3339 = time.RFC3339

// Policy bundles configurable security rules used by use cases.
type Policy struct {
	PasswordMinLength int
	AccessTTL         time.Duration
	RefreshTTL        time.Duration
	PasswordResetTTL  time.Duration
}

func (p Policy) ValidatePassword(plain string) error {
	if utf8.RuneCountInString(plain) < p.PasswordMinLength {
		return domainuser.ErrWeakPassword
	}
	return nil
}
