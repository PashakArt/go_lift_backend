SELECT
     ws.session_id,
     ws.started_at,
     ws.ended_at,
     ws.type,
     wset.set_id,
     wset.exercise_id,
     e.name,
     e.type,
     wset.set_number,
     wset.weight,
     wset.reps,
     wset.duration_seconds,
     wset.distance_meters
 FROM workout_sessions ws
 JOIN workout_sets wset ON wset.session_id = ws.session_id
 JOIN exercises e ON e.exercise_id = wset.exercise_id
 WHERE ws.is_active = false AND ws.user_id = $1
   AND ws.started_at >= $2 AND ws.started_at < $3
 ORDER BY ws.started_at, wset.exercise_id, wset.set_number;