INSERT INTO workout_sets (
    set_id,
    session_id,
    exercise_id,
    tenant_id,
    set_number,
    weight,
    reps,
    duration_seconds,
    distance_meters
) VALUES (
    $1,
    $2,
    $3,
    $4,
    COALESCE((
        SELECT MAX(ws.set_number) 
        FROM workout_sets ws 
        WHERE ws.session_id = $2 AND ws.exercise_id = $3
    ), 0) + 1,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (set_id) DO UPDATE SET
    weight           = EXCLUDED.weight,
    reps             = EXCLUDED.reps,
    duration_seconds = EXCLUDED.duration_seconds,
    distance_meters  = EXCLUDED.distance_meters
RETURNING 
    set_id, 
    session_id, 
    exercise_id, 
    tenant_id,
    set_number, 
    weight, 
    reps, 
    duration_seconds, 
    distance_meters, 
    created_at;