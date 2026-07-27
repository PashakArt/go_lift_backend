UPDATE workout_templates
SET name = $1,
    items = $2
WHERE template_id = $3 AND user_id = $4;