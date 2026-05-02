package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/lindritprekaj/notification-service/internal/domain/user"
)

//go:embed schema.sql
var schemaSQL string

type UserRepo struct {
	db *sql.DB
}

// Open opens (and creates if needed) the SQLite database at dbPath and runs
// the schema migration. The directory is created if it does not exist.
func Open(dbPath string) (*UserRepo, error) {
	if dbPath == "" {
		return nil, errors.New("sqlite: empty db path")
	}
	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		if err := ensureDir(dir); err != nil {
			return nil, fmt.Errorf("sqlite: ensure dir: %w", err)
		}
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(on)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	db.SetMaxOpenConns(1) // sqlite handles concurrency itself; one writer is plenty for our load

	if _, err := db.ExecContext(context.Background(), schemaSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}
	return &UserRepo{db: db}, nil
}

func (r *UserRepo) Close() error { return r.db.Close() }

func (r *UserRepo) Upsert(ctx context.Context, u user.User) error {
	if u.ID == "" || u.Email == "" {
		return errors.New("upsert: empty id or email")
	}
	ts := u.UpdatedAt
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users(user_id, email, updated_at) VALUES(?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET email = excluded.email, updated_at = excluded.updated_at
	`, u.ID, u.Email, ts)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetEmail(ctx context.Context, userID string) (string, error) {
	var email string
	err := r.db.QueryRowContext(ctx, `SELECT email FROM users WHERE user_id = ?`, userID).Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", user.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get email: %w", err)
	}
	return email, nil
}
