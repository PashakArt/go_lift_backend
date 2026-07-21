UPDATE workout_sessions
SET is_active = false, ended_at = NOW()
WHERE user_id = $1 and is_active = true