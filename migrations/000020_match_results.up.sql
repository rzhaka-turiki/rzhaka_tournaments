CREATE TABLE match_results (
    match_id UUID PRIMARY KEY REFERENCES matches(id) ON DELETE CASCADE,

    external_mid TEXT NOT NULL UNIQUE,
    map_name TEXT NOT NULL,
    map_id UUID REFERENCES maps(id) ON DELETE SET NULL,

    started_at TIMESTAMPTZ NOT NULL,

    aim_assist_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_match_results_match_id
    ON match_results(match_id);

CREATE INDEX idx_match_results_map_id
    ON match_results(map_id);

CREATE INDEX idx_match_results_started_at
    ON match_results(started_at);