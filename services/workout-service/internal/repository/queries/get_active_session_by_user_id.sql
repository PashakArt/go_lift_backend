SELECT session_id, tenant_id, user_id, template_id, type, started_at, ended_at, is_active
FROM workout_sessions
WHERE user_id = $1 AND is_active = true
