-- +goose Up
ALTER TABLE user_urine_metric
ADD COLUMN quantity NUMERIC(5, 2);

-- +goose Down
ALTER TABLE user_urine_metric
DROP COLUMN quantity;
