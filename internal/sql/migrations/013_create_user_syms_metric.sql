-- +goose Up
CREATE TABLE user_symptoms_metric (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    symptoms_id bigint NOT NULL REFERENCES symptoms ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    morning_severity NUMERIC(3,2) DEFAULT 0 CHECK (morning_severity >= 0 AND morning_severity <= 1),
    afternoon_severity NUMERIC(3,2) DEFAULT 0 CHECK (afternoon_severity >= 0 AND afternoon_severity <= 1),
    night_severity NUMERIC(3,2) DEFAULT 0 CHECK (night_severity >= 0 AND night_severity <= 1),
    CONSTRAINT unique_user_symptom_date UNIQUE (user_id, symptoms_id, date)
);

-- +goose Down
DROP TABLE IF EXISTS user_symptoms_metric;