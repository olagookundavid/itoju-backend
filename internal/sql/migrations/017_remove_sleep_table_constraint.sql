-- +goose Up
ALTER TABLE user_sleep_metric
DROP CONSTRAINT unique_user_sleep_date;


-- +goose Down
ALTER TABLE user_sleep_metric
ADD CONSTRAINT unique_user_sleep_date UNIQUE (user_id, is_night, date);