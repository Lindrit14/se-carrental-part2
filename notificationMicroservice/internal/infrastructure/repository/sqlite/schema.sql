CREATE TABLE IF NOT EXISTS users (
    user_id    TEXT PRIMARY KEY,
    email      TEXT NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
