-- Invite tokens table
CREATE TABLE IF NOT EXISTS invite_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    usage_count INTEGER NOT NULL DEFAULT 0,
    max_usage INTEGER NOT NULL DEFAULT 20,
    FOREIGN KEY (created_by) REFERENCES users(id)
);



