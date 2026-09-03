CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    map_id UUID NOT NULL REFERENCES maps(id),
    stats_token_id UUID REFERENCES match_api_tokens(id),

    group_id UUID REFERENCES groups(id),
    status TEXT NOT NULL DEFAULT 'pending',

    start_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_matches_map_id
    ON matches(map_id);

CREATE INDEX idx_matches_stats_token_id
    ON matches(stats_token_id);

CREATE INDEX idx_matches_start_at
    ON matches(start_at);

CREATE TABLE match_settings (
    match_id UUID PRIMARY KEY REFERENCES matches(id) ON DELETE CASCADE,

    drop_spots_enabled BOOLEAN NOT NULL DEFAULT FALSE,

    playlist_name TEXT NOT NULL,
    map_id UUID NOT NULL REFERENCES maps(id),
    admin_chat BOOLEAN NOT NULL DEFAULT FALSE,
    team_rename BOOLEAN NOT NULL DEFAULT FALSE,
    self_assign BOOLEAN NOT NULL DEFAULT TRUE,
    aim_assist BOOLEAN NOT NULL DEFAULT FALSE,
    anon_mode BOOLEAN NOT NULL DEFAULT FALSE,
    fill_bots_mode BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);