-- +goose Up
CREATE TABLE IF NOT EXISTS profile_pics ( 
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users ON DELETE CASCADE,
    url TEXT NOT NULL
);
     
-- +goose Down
DROP TABLE IF EXISTS profile_pics;