CREATE TABLE legends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    in_game_name TEXT NOT NULL UNIQUE,
    image_url TEXT,
    profile_image_url TEXT,
    class TEXT,
    ability TEXT,
    ultimate TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE match_legend_bans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    legend_id UUID NOT NULL REFERENCES legends(id),

    UNIQUE(match_id, legend_id)
);

CREATE INDEX idx_match_legend_bans_match_id
    ON match_legends_ban(match_id);