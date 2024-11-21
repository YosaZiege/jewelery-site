CREATE TABLE blacklisted_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL, -- Store the blacklisted token
    blacklisted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
