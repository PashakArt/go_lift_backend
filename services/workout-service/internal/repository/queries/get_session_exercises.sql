SELECT
     e.exercise_id,
     e.name,
     e.type,
     wset.set_id,
     wset.set_number,
     wset.weight,
     wset.reps,
     wset.duration_seconds,
     wset.distance_meters
 FROM workout_sets wset
 JOIN exercises e ON e.exercise_id = wset.exercise_id
 WHERE wset.session_id = $1
 ORDER BY e.name ASC, wset.set_number ASC;
