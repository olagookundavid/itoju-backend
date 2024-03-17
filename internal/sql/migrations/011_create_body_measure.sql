-- +goose Up
CREATE TABLE IF NOT EXISTS bodymeasure ( 
    user_id UUID NOT NULL UNIQUE REFERENCES users ON DELETE CASCADE,
    height SMALLINT NOT NULL DEFAULT 0,
    weight SMALLINT NOT NULL DEFAULT 0
);
     
-- +goose Down
DROP TABLE IF EXISTS bodymeasure;