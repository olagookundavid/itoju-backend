-- +goose Up
CREATE TABLE IF NOT EXISTS trackedmetrics ( 
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    CONSTRAINT unique_id_name UNIQUE (id, name)
     );
     
INSERT INTO trackedmetrics (name)
VALUES
    ('Food Diary'),
    ('Symptoms'),
    ('Sleep'),
    ('Menstruation and Ovulation'),
    ('Bowel Movements'),
    ('Medications'),
    ('Urination'),
    ('Exercise');
-- +goose Down
DROP TABLE IF EXISTS trackedmetrics;

