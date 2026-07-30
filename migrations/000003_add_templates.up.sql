ALTER TABLE workout_templates 
    RENAME COLUMN creator_id TO user_id;

ALTER TABLE workout_templates 
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS description;

ALTER TABLE workout_templates 
    ADD COLUMN IF NOT EXISTS items JSONB NOT NULL DEFAULT '[]'::jsonb;

CREATE INDEX IF NOT EXISTS idx_workout_templates_user_id 
    ON workout_templates(user_id);
