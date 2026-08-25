CREATE TABLE match_legend_bans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    legend_id UUID NOT NULL REFERENCES legends(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (match_id, legend_id)
);

CREATE INDEX idx_match_legend_bans_match_id
    ON match_legend_bans(match_id);

CREATE INDEX idx_match_legend_bans_legend_id
    ON match_legend_bans(legend_id);