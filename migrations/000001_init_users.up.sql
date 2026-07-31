CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    discord_id VARCHAR(32) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    avatar_url TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);