package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/lindritprekaj/notification-service/internal/domain/user"
)

func TestUserRepoUpsertAndGet(t *testing.T) {
	mr := miniredis.RunT(t)
	repo, err := New(mr.Addr(), "", 0, "notif")
	if err != nil {
		t.Fatalf("New: %v", err)
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

func TestUserRepoKeyPrefixIsolation(t *testing.T) {
	mr := miniredis.RunT(t)
	a, err := New(mr.Addr(), "", 0, "svcA")
	if err != nil {
		t.Fatalf("New svcA: %v", err)
	}
	defer a.Close()
	b, err := New(mr.Addr(), "", 0, "svcB")
	if err != nil {
		t.Fatalf("New svcB: %v", err)
	}
	defer b.Close()

	ctx := context.Background()
	if err := a.Upsert(ctx, user.User{ID: "u1", Email: "a@example.com"}); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if _, err := b.GetEmail(ctx, "u1"); !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("svcB must not see svcA's keys; got err=%v", err)
	}
}
