SELECT template_id, user_id, name, items, created_at
FROM workout_templates
WHERE user_id = $1
ORDER BY created_at DESC;