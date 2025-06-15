-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    username TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    type TEXT NOT NULL DEFAULT 'text',
    FOREIGN KEY (user_id) REFERENCES users(id)
);