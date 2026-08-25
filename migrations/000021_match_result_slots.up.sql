CREATE TABLE match_result_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    match_result_id UUID NOT NULL
        REFERENCES match_results(id) ON DELETE CASCADE,
    
    slot_number INT NOT NULL,

    team_name TEXT,
    team_placement INT,
    points NUMERIC(10, 2) NOT NULL DEFAULT 0,
    kills INT NOT NULL DEFAULT 0,

    match_slot_id UUID
        REFERENCES match_slots(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (match_result_id, slot_number)
);

CREATE INDEX idx_match_result_slots_result_id
    ON match_result_slots(match_result_id);

CREATE INDEX idx_match_result_slots_match_slot_id
    ON match_result_slots(match_slot_id);
