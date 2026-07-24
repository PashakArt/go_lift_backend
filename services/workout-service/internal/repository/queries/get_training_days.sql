SELECT DISTINCT TO_CHAR(DATE(started_at), 'YYYY-MM-DD') as workout_date
FROM workout_sessions
WHERE user_id = $1 
  AND is_active = false
  AND started_at >= $2 AND started_at < $3;