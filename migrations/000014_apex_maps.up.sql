CREATE TABLE maps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    in_game_name TEXT NOT NULL UNIQUE,
    image_url TEXT,
    minimap_image_url TEXT,
    supports_drop_spots BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE map_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    map_id UUID NOT NULL REFERENCES maps(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    image_url TEXT,
    position INT NOT NULL,
    x DOUBLE PRECISION NOT NULL,
    y DOUBLE PRECISION NOT NULL,

    UNIQUE (map_id, name),
    UNIQUE (map_id, position)
);

CREATE INDEX idx_map_locations_map_id
    ON map_locations(map_id);