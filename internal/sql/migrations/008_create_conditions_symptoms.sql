-- +goose Up
CREATE TABLE IF NOT EXISTS user_symptoms (
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    symptoms_id bigint NOT NULL REFERENCES symptoms ON DELETE CASCADE,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, symptoms_id)
);

CREATE TABLE IF NOT EXISTS user_conditions (
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    conditions_id bigint NOT NULL REFERENCES conditions ON DELETE CASCADE,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, conditions_id)
);


-- +goose Down
DROP TABLE IF EXISTS user_symptoms;
DROP TABLE IF EXISTS user_conditions;