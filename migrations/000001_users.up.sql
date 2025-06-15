-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    chat_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    first_name VARCHAR(255),
    last_name VARCHAR(255)
);


