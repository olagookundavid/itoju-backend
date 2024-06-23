-- +goose Up
ALTER TABLE user_urine_metric
    ALTER COLUMN pain TYPE NUMERIC(3,2) USING pain::NUMERIC(3,2),
    ALTER COLUMN pain SET DEFAULT 0,
    ADD CONSTRAINT urine_pain_check CHECK (pain >= 0 AND pain <= 1);

-- +goose Down
ALTER TABLE user_urine_metric
    ALTER COLUMN pain TYPE SMALLINT USING pain::smallint,
    ALTER COLUMN pain SET DEFAULT 0,
    DROP CONSTRAINT urine_pain_check;
