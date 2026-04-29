package user

import "context"

// Repository is the persistence contract for User.
// Infrastructure adapters implement it; the application layer depends on it.
type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error

	// List returns a page of users sorted by created_at desc, plus the
	// total document count.
	List(ctx context.Context, limit, offset int) (items []*User, total int64, err error)

	// CountByRole returns how many users currently hold the given role.
	CountByRole(ctx context.Context, role Role) (int64, error)
}
