CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL, -- When the refresh token expires
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
