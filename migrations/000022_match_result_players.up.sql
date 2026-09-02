CREATE TABLE match_result_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_result_slot_id UUID NOT NULL
        REFERENCES match_result_slots(id) ON DELETE CASCADE,
    
    nid_hash TEXT NOT NULL,
    player_name TEXT NOT NULL,

    legend_id UUID
        REFERENCES legends(id) ON DELETE SET NULL,

    character_name TEXT NOT NULL,

    kills INT NOT NULL DEFAULT 0,
    assists INT NOT NULL DEFAULT 0,
    knockdowns INT NOT NULL DEFAULT 0,
    damage_dealt INT NOT NULL DEFAULT 0,
    survival_time INT NOT NULL DEFAULT 0,
    hardware TEXT NOT NULL DEFAULT NULL,
    headshots INT NOT NULL DEFAULT 0,
    shots INT NOT NULL DEFAULT 0,
    hits INT NOT NULL DEFAULT 0,
    respawns INT NOT NULL DEFAULT 0,
    revives INT NOT NULL DEFAULT 0,


    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_match_result_players_slot_id
    ON match_result_players(match_result_slot_id);

CREATE INDEX idx_match_result_players_nid_hash
    ON match_result_players(nid_hash);

CREATE INDEX idx_match_result_players_legend_id
    ON match_result_players(legend_id);