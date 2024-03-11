-- +goose Up
CREATE TABLE IF NOT EXISTS symptoms ( 
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
     );
     
INSERT INTO symptoms (name)
VALUES
        ('Abdominal Pain (Left)'),
		('Abdominal Pain (Lower)'),
		('Abdominal Pain (Right)'),
		('Abdominal Pain (Upper)'),
		('Acne'),
		('Back Pain (Lower)'),
		('Back Pain (Upper)'),
		('Bleeding'),
		('Bloating'),
		('Chest Pain'),
		('Cold & Catarrh'),
		('Constipation'),
		('Cough'),
		('Diarrhea'),
		('Earache'),
		('Eczema'),
		('Eye Pain'),
		('Fatigue'),
		('Fever'),
		('Flu'),
		('Headache'),
		('Indigestion'),
		('Insomnia'),
		('Joint Pain'),
		('Knee Pain'),
		('Leg Pain'),
		('Loss of Smell'),
		('Loss of Taste'),
		('Menstrual Pain (Just Before & During)'),
		('Migraine'),
		('Mouth Sores'),
		('Muscle Pain'),
		('Nausea'),
		('Neck Pain'),
		('Pain During Bowel Movements'),
		('Pain During Sex'),
		('Pain During Urination'),
		('Piles'),
		('Rectal bleeding'),
		('Sciatica'),
		('Seizure'),
		('Shortness of Breath'),
		('Shoulder Pain (Left)'),
		('Shoulder Pain (Right)'),
		('Sore Throat'),
		('Stomach Cramps'),
		('Tinnitus'),
		('Toothache'),
		('Under-rib Pain'),
		('Vertigo'),
		('Vomiting');

-- +goose Down
DROP TABLE IF EXISTS symptoms;

