-- +goose Up
CREATE TABLE cycles_days (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    cycle_id UUID NOT NULL REFERENCES menstrual_cycles ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    date DATE NOT NULL,
    is_period BOOLEAN NOT NULL,
    is_ovulation BOOLEAN NOT NULL,
    flow NUMERIC(3,2) NOT NULL,
    pain  NUMERIC(3,2) NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    cmq TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS cycles_days;