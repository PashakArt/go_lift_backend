SELECT
    template_id,
    user_id,
    name,
    items,
    created_at
FROM workout_templates
WHERE template_id = $1 AND user_id = $2;