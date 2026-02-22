
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description  TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);