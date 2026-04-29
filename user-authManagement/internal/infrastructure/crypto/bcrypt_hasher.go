// Package crypto provides infrastructure adapters for password hashing
// and JWT signing/parsing.
package crypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{ Cost int }

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = 12
	}
	return &BcryptHasher{Cost: cost}
}

func (b *BcryptHasher) Hash(plain string) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(plain), b.Cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(out), nil
}

func (b *BcryptHasher) Verify(plain, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return err
		}
		return fmt.Errorf("bcrypt verify: %w", err)
	}
	return nil
}
