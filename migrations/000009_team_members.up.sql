CREATE TABLE team_members (
    team_id UUID NOT NULL
        REFERENCES teams(id)
        ON DELETE CASCADE,
    
    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    role TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (team_id, user_id)
);

CREATE INDEX team_members_user_id_idx
ON team_members(user_id);