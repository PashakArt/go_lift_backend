CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    user_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    telegram_id VARCHAR(64) NOT NULL,
    phone       VARCHAR(32),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_tenant_telegram UNIQUE (tenant_id, telegram_id)
);

CREATE TABLE IF NOT EXISTS tenants (
    tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_token_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    branding_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_templates (
    template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    creator_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TYPE session_type AS ENUM ('classic', 'circuit');

CREATE TABLE IF NOT EXISTS workout_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    template_id UUID REFERENCES workout_templates(template_id) ON DELETE SET NULL,
    type session_type NOT NULL DEFAULT 'classic',
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN default true
);

CREATE TABLE IF NOT EXISTS muscle_groups (
    muscle_group_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL UNIQUE, -- 'chest', 'triceps', 'biceps'
    name VARCHAR(64) NOT NULL,     -- 'Грудь', 'Трицепс', 'Бицепс'
    sort_order INT NOT NULL DEFAULT 0
);

INSERT INTO muscle_groups (code, name, sort_order) VALUES
    ('shoulders', 'Плечи', 10),
    ('chest', 'Грудь', 20),
    ('back', 'Спина', 30),
    ('biceps', 'Бицепс', 40),
    ('triceps', 'Трицепс', 50),
    ('forearms', 'Предплечья', 60),
    ('abs', 'Пресс', 70),
    ('glutes', 'Ягодицы', 80),
    ('quads', 'Квадрицепсы (Передняя поверхность бедра)', 90),
    ('hamstrings', 'Бицепс бедра (Задняя поверхность бедра)', 100),
    ('adductors', 'Приводящие мышцы бедра (Внутренняя поверхность)', 110),
    ('calves', 'Икры', 120)
ON CONFLICT (code) DO UPDATE 
SET name = EXCLUDED.name,
    sort_order = EXCLUDED.sort_order;

CREATE TYPE exercise_type AS ENUM ('dynamic', 'static', 'bodyweight', 'cardio');

CREATE TABLE IF NOT EXISTS exercises (
    exercise_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type exercise_type NOT NULL,
    is_global BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS exercise_muscle_groups (
    exercise_id UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    muscle_group_id UUID REFERENCES muscle_groups(muscle_group_id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, muscle_group_id)
);


CREATE TABLE IF NOT EXISTS exercise_muscle_groups (
    exercise_id UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    muscle_group_id UUID REFERENCES muscle_groups(muscle_group_id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, muscle_group_id)
);

CREATE TABLE IF NOT EXISTS user_favorite_exercises (
    user_id     UUID REFERENCES users(user_id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    
    PRIMARY KEY (user_id, exercise_id)
);

CREATE INDEX IF NOT EXISTS idx_user_favorites_user_id ON user_favorite_exercises(user_id);

-- 1. Создаем временную функцию или блок, чтобы удобно связать по текстовым кодам (чистый SQL)
DO $$
DECLARE
    -- Переменные для хранения ID упражнений
    bench_press_id UUID;
    pushups_id UUID;
    pullups_id UUID;
    barbell_row_id UUID;
    squats_id UUID;
    leg_press_id UUID;
    bicep_curl_id UUID;
    tricep_extension_id UUID;
    plank_id UUID;
    adductor_machine_id UUID;
BEGIN

    -- ==========================================
    -- ГРУДЬ & ТРИЦЕПС
    -- ==========================================
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Жим штанги лежа', 'dynamic', true, NULL) RETURNING exercise_id INTO bench_press_id;
    
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Отжимания от пола', 'bodyweight', true, NULL) RETURNING exercise_id INTO pushups_id;

    -- Связываем Жим лежа с Грудью (основная) и Трицепсом (синергист)
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (bench_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'chest')),
        (bench_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps'));

    -- Связываем Отжимания с Грудью и Трицепсом
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (pushups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'chest')),
        (pushups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps'));


    -- ==========================================
    -- СПИНА & БИЦЕПС
    -- ==========================================
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Подтягивания на турнике', 'bodyweight', true, NULL) RETURNING exercise_id INTO pullups_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Тяга штанги в наклоне', 'dynamic', true, NULL) RETURNING exercise_id INTO barbell_row_id;

    -- Связываем Подтягивания со Спиной и Бицепсом
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (pullups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'back')),
        (pullups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'biceps'));

    -- Связываем Тягу штанги со Спиной
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (barbell_row_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'back'));


    -- ==========================================
    -- НОГИ & ПРИВОДЯЩИЕ
    -- ==========================================
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Приседания со штангой', 'dynamic', true, NULL) RETURNING exercise_id INTO squats_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Жим ногами в тренажере', 'dynamic', true, NULL) RETURNING exercise_id INTO leg_press_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Сведение ног в тренажере', 'dynamic', true, NULL) RETURNING exercise_id INTO adductor_machine_id;

    -- Приседания: Квадрицепсы, Ягодицы, Бицепс бедра, Приводящие (работают во всех тяжелых седах!)
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'quads')),
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes')),
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'hamstrings')),
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'adductors'));

    -- Жим ногами: Квадрицепсы, Ягодицы
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (leg_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'quads')),
        (leg_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes'));

    -- Сведение ног: строго Приводящие мышцы
    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (adductor_machine_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'adductors'));


    -- ==========================================
    -- ИЗОЛЯЦИЯ: РУКИ И КОР
    -- ==========================================
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Подъем штанги на бицепс', 'dynamic', true, NULL) RETURNING exercise_id INTO bicep_curl_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Разгибания на трицепс на верхнем блоке', 'dynamic', true, NULL) RETURNING exercise_id INTO tricep_extension_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Планка классическая', 'static', true, NULL) RETURNING exercise_id INTO plank_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (bicep_curl_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'biceps'));

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (tricep_extension_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps'));

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (plank_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'abs')),
        (plank_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'core'));

END $$;