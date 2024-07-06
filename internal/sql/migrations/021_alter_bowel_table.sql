-- +goose Up
ALTER TABLE user_bowel_metric
    ALTER COLUMN pain TYPE NUMERIC(3,2) USING pain::NUMERIC(3,2),
    ALTER COLUMN pain SET DEFAULT 0,
    ADD CONSTRAINT pain_check CHECK (pain >= 0 AND pain <= 1);

-- +goose Down
ALTER TABLE user_bowel_metric
    ALTER COLUMN pain TYPE SMALLINT USING pain::smallint,
    ALTER COLUMN pain SET DEFAULT 0,
    DROP CONSTRAINT pain_check;
