
CREATE TABLE IF NOT EXISTS sessions (
    id              SERIAL PRIMARY KEY,
    device          VARCHAR(255),
    refresh_token   TEXT NOT NULL,
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address      VARCHAR(50),
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

