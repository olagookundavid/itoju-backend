-- +goose Up
CREATE TABLE IF NOT EXISTS menstruation ( 
    user_id UUID NOT NULL UNIQUE REFERENCES users ON DELETE CASCADE,
    period_len SMALLINT NOT NULL DEFAULT 0,
    cycle_len SMALLINT NOT NULL DEFAULT 0
);
     
-- +goose Down
DROP TABLE IF EXISTS menstruation;