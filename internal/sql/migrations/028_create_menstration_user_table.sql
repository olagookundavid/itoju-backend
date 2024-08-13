-- +goose Up
CREATE TABLE menstrual_cycles (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    start_date DATE NOT NULL,
    cycle_length SMALLINT DEFAULT 28,
    period_length SMALLINT DEFAULT 5,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS menstrual_cycles;