// Package ports defines the outbound interfaces (driven ports) that
// the application layer needs from the outside world.
// Infrastructure adapters implement these.
package ports

// PasswordHasher abstracts password hashing/verification.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hash string) error
}
