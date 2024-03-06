-- +goose Up
CREATE TABLE IF NOT EXISTS smiley ( 
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    CONSTRAINT unique_id_smileyname UNIQUE (id, name)
     );
     
INSERT INTO smiley (name)
VALUES
    ('Very Good'),
    ('Good'),
    ('Mild'),
    ('Bad'),
    ('Very Bad');

-- +goose Down
DROP TABLE IF EXISTS smiley;

