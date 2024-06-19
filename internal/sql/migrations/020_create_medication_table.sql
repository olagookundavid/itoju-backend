-- +goose Up
CREATE TABLE user_medication_metric (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    name TEXT NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    dosage SMALLINT NOT NULL,
    metric TEXT NOT NULL,
    quantity SMALLINT DEFAULT 0,
    time VARCHAR(20) NOT NULL DEFAULT ''
);


-- +goose Down
DROP TABLE IF EXISTS user_medication_metric;