-- +goose Up
CREATE TABLE IF NOT EXISTS conditions ( 
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
     );
     
INSERT INTO conditions (name)
VALUES
    ('Premenstrual syndrome (PMS)'),
    ('Ovarian cysts'),
    ('Endometriosis'),
    ('Uterine fibroids'),
    ('Polycystic ovarian syndrome (PCOS)'),
    ('Adenomyosis'),
    ('Urinary Incontinence'),
    ('Infertility'),
    ('Uterine Prolapse'),
    ('Cervical Cancer'),
    ('Ovarian Cancer'),
    ('Endometrial cancer'),
    ('Fibromyalgia'),
    ('Interstitial cystitis'),
    ('Dysmenorrhea (Painful periods)'),
    ('Amenorrhea (Absence of periods)'),
    ('Vaginismus'),
    ('Pelvic Inflammatory Disease'),
    ('Vulvodynia'),
    ('Menopause'),
    ('Vaginitis');
 


-- +goose Down
DROP TABLE IF EXISTS conditions;

