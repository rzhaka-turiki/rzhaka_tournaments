CREATE TABLE match_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,

    slot_number INT NOT NULL,
    drop_spot_id UUID REFERENCES map_locations(id),

    created_at TIMETAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (match_id, slot_number)
);

CREATE INDEX idx_match_slots_match_id
    ON match_slots(match_id);

CREATE INDEX idx_match_slots_drop_spot_id
    ON match_slots(drop_spot_id);

CREATE TABLE match_slot_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_slot_id UUID NOT NULL REFERENCES match_slots(id) ON DELETE CASCADE,

    expected_nid_hash TEXT,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updaeted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (match_slot_id, user_id)
);

CREATE INDEX idx_match_slot_players_slot_id
    ON match_slot_players(match_slot_id);

CREATE INDEX idx_match_slot_players_expected_nid_hash
    ON match_slot_players(expected_nid_hash);

CREATE INDEX idx_match_slot_players_user_id
    ON match_slot_players(user_id);