-- +goose Up
CREATE TABLE user_sleep_metric (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    is_night BOOLEAN NOT NULL,
    time_slept TIME,
    time_woke_up TIME, 
    tags TEXT[] NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    severity NUMERIC(3,2) DEFAULT 0 CHECK (severity >= 0 AND severity <= 1),
    CONSTRAINT unique_user_sleep_date UNIQUE (user_id, is_night, date)
);


-- +goose Down
DROP TABLE IF EXISTS user_sleep_metric;