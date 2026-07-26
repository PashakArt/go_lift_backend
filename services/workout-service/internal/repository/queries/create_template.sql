INSERT INTO workout_templates (user_id, name, items)
VALUES ($1, $2, $3)
RETURNING template_id;