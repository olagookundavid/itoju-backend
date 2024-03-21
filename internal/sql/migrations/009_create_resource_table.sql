-- +goose Up
CREATE TABLE IF NOT EXISTS resources ( 
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    imageUrl TEXT NOT NULL,
    link TEXT NOT NULL ,
    tags TEXT[] NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 
);

-- +goose Down
DROP TABLE IF EXISTS resources;