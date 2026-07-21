WITH inserted AS (
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
    ON CONFLICT (user_id) WHERE is_active = true DO NOTHING
    RETURNING session_id, started_at
)
SELECT session_id, started_at FROM inserted
UNION ALL
SELECT session_id, started_at 
FROM workout_sessions
WHERE user_id = $2 AND is_active = true
LIMIT 1;