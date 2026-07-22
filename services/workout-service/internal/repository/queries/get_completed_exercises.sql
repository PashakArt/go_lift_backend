SELECT 
    ws.set_id,
    ws.set_number,
    ws.weight,
    ws.reps,
    ws.duration_seconds,
    ws.distance_meters
FROM workout_sets ws
JOIN workout_sessions s ON ws.session_id = s.session_id
WHERE s.user_id = $1 
  AND s.is_active = true
  AND ws.exercise_id = $2
ORDER BY ws.set_number ASC;