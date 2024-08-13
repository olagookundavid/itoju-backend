-- +goose Up
CREATE TABLE cycles_days (
    id SERIAL PRIMARY KEY,
    cycle_id INTEGER NOT NULL REFERENCES menstrual_cycles ON DELETE CASCADE,
    date DATE NOT NULL,
    is_period BOOLEAN NOT NULL,
    is_ovulation BOOLEAN NOT NULL,
    flow SMALLINT NOT NULL,
    pain SMALLINT NOT NULL,
    tags TEXT[] DEFAULT '{}',
    cmq TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS cycles_days;