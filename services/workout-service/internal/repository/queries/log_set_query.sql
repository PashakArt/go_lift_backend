INSERT INTO workout_sets (
    set_id,
    session_id,
    exercise_id,
    set_number,
    weight,
    reps,
    duration_seconds,
    distance_meters
) VALUES (
    COALESCE(NULLIF($1, ''), gen_random_uuid()),
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (set_id) DO UPDATE SET
    set_number       = EXCLUDED.set_number,
    weight           = EXCLUDED.weight,
    reps             = EXCLUDED.reps,
    duration_seconds = EXCLUDED.duration_seconds,
    distance_meters  = EXCLUDED.distance_meters
RETURNING 
    set_id, 
    session_id, 
    exercise_id, 
    set_number, 
    weight, 
    reps, 
    duration_seconds, 
    distance_meters, 
    created_at;