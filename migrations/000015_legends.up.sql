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