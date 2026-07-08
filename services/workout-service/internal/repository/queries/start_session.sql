INSERT INTO workout_sessions (
    tenant_id,
    user_id,
    template_id,
    type,
    started_at
) VALUES (
    $1,
    $2,
    NULLIF($3, '')::uuid,
    $4,
    NOW()
)
RETURNING session_id, started_at;