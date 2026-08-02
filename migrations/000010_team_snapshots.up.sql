CREATE TABLE team_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    team_id UUID NOT NULL
        REFERENCES teams(id)
        ON DELETE CASCADE,
    
    name TEXT NOT NULL,
    short_name VARCHAR(5),
    logo_path TEXT,
    logo_dark_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);