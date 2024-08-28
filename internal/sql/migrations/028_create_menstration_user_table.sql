-- +goose Up
CREATE TABLE menstrual_cycles (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    start_date DATE UNIQUE NOT NULL,
    cycle_length SMALLINT DEFAULT 28,
    period_length SMALLINT DEFAULT 5,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS menstrual_cycles;