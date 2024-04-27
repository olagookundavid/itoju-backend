-- +goose Up
CREATE TABLE user_exercise_metric (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    name TEXT NOT NULL,
    started VARCHAR(20) NOT NULL DEFAULT '',
    ended VARCHAR(20) NOT NULL DEFAULT '', 
    tags TEXT[] NOT NULL DEFAULT '{}',
    no_of_times SMALLINT NOT NULL DEFAULT 0
);


-- +goose Down
DROP TABLE IF EXISTS user_exercise_metric;