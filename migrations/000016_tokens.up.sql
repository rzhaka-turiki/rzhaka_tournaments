CREATE TABLE match_api_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_api_token_id INT NOT NULL UNIQUE,
    added_by UUID NOT NULL REFERENCES users(id),
    organisation_id NOT NULL UUID REFERENCES organisations(id),

    activation TIMESTAMPTZ NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,

    stats_token TEXT NOT NULL,
    admin_token TEXT,
    player_token TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);