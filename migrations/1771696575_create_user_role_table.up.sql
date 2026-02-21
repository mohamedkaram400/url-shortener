
CREATE TABLE IF NOT EXISTS user_role (
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id         INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);