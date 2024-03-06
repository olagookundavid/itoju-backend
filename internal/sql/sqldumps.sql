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
