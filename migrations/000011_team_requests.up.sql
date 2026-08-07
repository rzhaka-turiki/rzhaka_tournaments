CREATE TABLE team_invitations (
    id UUID PRIMARY KEY,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (
        type IN (
            'invite',
            'invite_link'
        )
    ),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(team_id, user_id)
);

CREATE INDEX idx_team_requests_team_id
ON team_requests(team_id);

CREATE INDEX idx_team_requests_user_id
ON team_requests(user_id);

CREATE INDEX idx_team_requests_expires_at
ON team_requests(expires_at);