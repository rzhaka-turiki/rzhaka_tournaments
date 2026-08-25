CREATE TABLE match_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id UUID REFERENCES matches(id) ON DELETE SET NULL,

    external_mid TEXT NOT NULL UNIQUE,
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