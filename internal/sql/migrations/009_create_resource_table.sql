-- +goose Up
CREATE TABLE IF NOT EXISTS resources ( 
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    imageUrl TEXT NOT NULL UNIQUE,
    link TEXT NOT NULL UNIQUE,
    tags text[] NOT NULL
     );
     
-- +goose Down
DROP TABLE IF EXISTS resources;