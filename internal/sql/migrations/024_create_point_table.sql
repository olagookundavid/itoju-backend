-- +goose Up
CREATE TABLE user_point (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    tot_point bigint NOT NULL,
    created_date DATE DEFAULT CURRENT_DATE,
    CONSTRAINT unique_user_point UNIQUE (user_id)
);

-- +goose Down
DROP TABLE IF EXISTS user_point;