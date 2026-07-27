SELECT
    exercise_id,
    name,
    type,
    is_global,
    created_at
FROM exercises
WHERE exercise_id = ANY($1);