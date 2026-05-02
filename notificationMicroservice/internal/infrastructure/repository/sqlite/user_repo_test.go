package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/lindritprekaj/notification-service/internal/domain/user"
)

func TestUserRepoUpsertAndGet(t *testing.T) {
	dir := t.TempDir()
	repo, err := Open(filepath.Join(dir, "users.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	if _, err := repo.GetEmail(ctx, "missing"); !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	now := time.Now().UTC()
	if err := repo.Upsert(ctx, user.User{ID: "u1", Email: "a@example.com", UpdatedAt: now}); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	got, err := repo.GetEmail(ctx, "u1")
	if err != nil {
		t.Fatalf("GetEmail: %v", err)
	}
	if got != "a@example.com" {
		t.Fatalf("want a@example.com, got %s", got)
	}

	// Upsert overwrites email
	if err := repo.Upsert(ctx, user.User{ID: "u1", Email: "b@example.com", UpdatedAt: now}); err != nil {
		t.Fatalf("Upsert overwrite: %v", err)
	}
	got, err = repo.GetEmail(ctx, "u1")
	if err != nil {
		t.Fatalf("GetEmail: %v", err)
	}
	if got != "b@example.com" {
		t.Fatalf("want b@example.com, got %s", got)
	}
}
