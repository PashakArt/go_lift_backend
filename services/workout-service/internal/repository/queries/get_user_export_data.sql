SELECT 
    ws.session_id,
    ws.started_at,
    ws.ended_at,
    ws.type AS session_type,
    e.name AS exercise_name,
    st.set_number,
    st.weight,
    st.reps,
    st.duration_seconds,
    st.distance_meters
FROM workout_sessions ws
JOIN workout_sets st ON st.session_id = ws.session_id
JOIN exercises e ON e.exercise_id = st.exercise_id
WHERE ws.user_id = $1
ORDER BY ws.started_at DESC, st.created_at ASC