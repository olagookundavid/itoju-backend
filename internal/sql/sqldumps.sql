insert

INSERT INTO user_trackedmetric (user_id, metric_id)
SELECT 'user_id_here', id
FROM trackedmetrics;

get list
SELECT t.name
FROM trackedmetrics t
JOIN user_trackedmetric utm ON t.id = utm.metric_id
WHERE utm.user_id = 'user_id_here';

DELETE
DELETE FROM user_trackedmetric
WHERE user_id = 'user_id_here'
AND metric_id = metric_id_to_delete;


LIST OF SYMPTOMS

Abdominal Pain (Left)
Abdominal Pain (Lower)
Abdominal Pain (Right)
Abdominal Pain (Upper)
Acne
Back Pain (Lower)
Back Pain (Upper)
Bleeding
Bloating
Chest Pain
Cold & Catarrh
Constipation
Cough
Diarrhea
Earache
Eczema
Eye Pain
Fatigue
Fever
Flu
Headache
Indigestion
Insomnia
Joint Pain
Knee Pain
Leg Pain
Loss of Smell
Loss of Taste
Menstrual Pain (Just Before & During)
Migraine
Mouth Sores
Muscle Pain
Nausea
Neck Pain
Pain During Bowel Movements
Pain During Sex
Pain During Urination 
Piles
Rectal bleeding
Sciatica
Seizure
Shortness of Breath
Shoulder Pain (Left)
Shoulder Pain (Right)
Sore Throat
Stomach Cramps
Tinnitus
Toothache
Under-rib Pain
Vertigo
Vomiting



//
LIST OF CONDITIONS

- Pre-menstrual syndrome (PMS)
- Ovarian cysts
- Endometriosis
- Uterine fibroids
- Polycystic ovarian syndrome (PCOS)
- Adenomyosis
- Urinary Incontinence
- Infertility
- Uterine Prolapse
- Cervical Cancer
- Ovarian Cancer
- Endometrial cancer
- Fibromyalgia
- Interstitial cystitis 
- Dysmenorrhea (Painful periods)
- Amenorrhea (Absence of periods)
- Vaginitis
- Vaginismus
- Pelvic Inflammatory Disease 
- Vulvodynia
- Menopause