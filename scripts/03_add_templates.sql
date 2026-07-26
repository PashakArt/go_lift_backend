ALTER TABLE workout_templates 
    RENAME COLUMN creator_id TO user_id;

ALTER TABLE workout_templates 
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS description;

ALTER TABLE workout_templates 
    ADD COLUMN IF NOT EXISTS items JSONB NOT NULL DEFAULT '[]'::jsonb;

CREATE INDEX IF NOT EXISTS idx_workout_templates_user_id 
    ON workout_templates(user_id);

-- [
--   {
--     "exercise_id": "550e8400-e29b-41d4-a716-446655440000",
--     "order_index": 1,
--     "target_sets": [
--       {
--         "set_num": 1,
--         "weight": 80.0,
--         "reps": 10
--       },
--       {
--         "set_num": 2,
--         "weight": 85.0,
--         "reps": 8
--       }
--     ]
--   }
-- ]