-- +goose Up
CREATE TABLE user_food_metric (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    breakfast_meal TEXT NOT NULL DEFAULT '',
    lunch_meal TEXT NOT NULL DEFAULT '',
    dinner_meal TEXT NOT NULL DEFAULT '',
    breakfast_extra TEXT NOT NULL DEFAULT '',
    lunch_extra TEXT NOT NULL DEFAULT '',
    dinner_extra TEXT NOT NULL DEFAULT '',
    breakfast_fruit TEXT NOT NULL DEFAULT '',
    lunch_fruit TEXT NOT NULL DEFAULT '',
    dinner_fruit TEXT NOT NULL DEFAULT '',
    breakfast_tags TEXT[] NOT NULL DEFAULT '{}',
    lunch_tags TEXT[] NOT NULL DEFAULT '{}',
    dinner_tags TEXT[] NOT NULL DEFAULT '{}',
    snack_name TEXT NOT NULL DEFAULT '',
    snack_tags TEXT[] NOT NULL DEFAULT '{}',
    glass_no SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT unique_user_food_date UNIQUE (user_id, date)
);


-- +goose Down
DROP TABLE IF EXISTS user_food_metric;