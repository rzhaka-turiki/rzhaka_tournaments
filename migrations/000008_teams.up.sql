CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    short_name VARCHAR(5),

    logo_path TEXT,
    logo_dark_path TEXT,

    owner_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    archived_at TIMESTAMPTZ
);
CREATE INDEX idx_teams_owner
ON teams(owner_id);