-- +goose Up
CREATE TABLE user_urine_metric (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    type SMALLINT NOT NULL,
    pain SMALLINT DEFAULT 0,
    time VARCHAR(20) NOT NULL DEFAULT '',
    tags TEXT[] NOT NULL DEFAULT '{}'
);


-- +goose Down
DROP TABLE IF EXISTS user_urine_metric;