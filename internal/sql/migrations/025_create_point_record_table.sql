-- +goose Up
CREATE TABLE user_point_record (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    point bigint NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE
);

-- +goose Down
DROP TABLE IF EXISTS user_point_record;