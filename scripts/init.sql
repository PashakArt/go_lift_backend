CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- 1. СОЗДАНИЕ ТАБЛИЦ
-- ==========================================

CREATE TABLE IF NOT EXISTS tenants (
    tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    branding_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE tenants ADD CONSTRAINT unique_tenant_code UNIQUE (code);

CREATE TABLE IF NOT EXISTS users (
    user_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    telegram_id VARCHAR(64) NOT NULL,
    phone       VARCHAR(32),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_tenant_telegram UNIQUE (tenant_id, telegram_id)
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
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(64) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0
);

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

CREATE TABLE IF NOT EXISTS user_favorite_exercises (
    user_id     UUID REFERENCES users(user_id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, exercise_id)
);

CREATE INDEX IF NOT EXISTS idx_user_favorites_user_id ON user_favorite_exercises(user_id);


-- ==========================================
-- 2. НАПОЛНЕНИЕ ДЕФОЛТНЫМИ ДАННЫМИ
-- ==========================================

-- Дефолтный тенант
INSERT INTO tenants (tenant_id, name, code, branding_json)
VALUES (
    '00000000-0000-0000-0000-000000000000', 
    'GoLift', 
    'default', 
    '{
        "theme": {
            "mode": "dark",
            "primary_color": "#FF5722",
            "background_color": "#121212",
            "surface_color": "#1E1E1E",
            "text_color": "#FFFFFF",
            "accent_color": "#00E676"
        },
        "assets": {
            "logo_url": "https://assets.golift.app/logos/default_logo.png",
            "welcome_image_url": "https://assets.golift.app/images/default_welcome.png"
        }
    }'::jsonb
)
ON CONFLICT (tenant_id) DO NOTHING;

-- Справочник мышечных групп
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


-- Аннонимный блок для заливки ВСЕХ упражнений
DO $$
DECLARE
    def_tenant_id UUID := '00000000-0000-0000-0000-000000000000';
    
    -- Переменные под ID упражнений
    bench_press_id UUID;
    pushups_id UUID;
    pullups_id UUID;
    barbell_row_id UUID;
    squats_id UUID;
    leg_press_id UUID;
    adductor_machine_id UUID;
    bicep_curl_id UUID;
    tricep_extension_id UUID;
    plank_id UUID;
    overhead_press_id UUID;
    lateral_raise_id UUID;
    lat_pulldown_id UUID;
    hyperextension_id UUID;
    hip_thrust_id UUID;
    lunges_id UUID;
    calf_raise_id UUID;
    crunch_id UUID;
    hanging_legs_id UUID;
    hammer_curl_id UUID;
BEGIN

    -- 1. ГРУДЬ И ТРИЦЕПС
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Жим штанги лежа', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO bench_press_id;
    
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Отжимания от пола', 'bodyweight', true, def_tenant_id) RETURNING exercise_id INTO pushups_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (bench_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'chest')),
        (bench_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps')),
        (pushups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'chest')),
        (pushups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps'));

    -- 2. СПИНА И БИЦЕПС
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Подтягивания на турнике', 'bodyweight', true, def_tenant_id) RETURNING exercise_id INTO pullups_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Тяга штанги в наклоне', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO barbell_row_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Тяга верхнего блока к груди', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO lat_pulldown_id;
    
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Гиперэкстензия', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO hyperextension_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (pullups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'back')),
        (pullups_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'biceps')),
        (barbell_row_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'back')),
        (lat_pulldown_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'back')),
        (lat_pulldown_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'biceps')),
        (hyperextension_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'back')),
        (hyperextension_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes'));

    -- 3. НОГИ, ЯГОДИЦЫ И ПРИВОДЯЩИЕ
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Приседания со штангой', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO squats_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Жим ногами в тренажере', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO leg_press_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Сведение ног в тренажере', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO adductor_machine_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Ягодичный мостик со штангой', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO hip_thrust_id;
    
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Выпады с гантелями', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO lunges_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'quads')),
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes')),
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'hamstrings')),
        (squats_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'adductors')),
        (leg_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'quads')),
        (leg_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes')),
        (adductor_machine_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'adductors')),
        (hip_thrust_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes')),
        (hip_thrust_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'hamstrings')),
        (lunges_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'quads')),
        (lunges_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes'));

    -- 4. ПЛЕЧИ
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Армейский жим (Жим штанги стоя)', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO overhead_press_id;
    
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Махи гантелями в стороны', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO lateral_raise_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (overhead_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'shoulders')),
        (overhead_press_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps')),
        (lateral_raise_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'shoulders'));

    -- 5. ИЗОЛЯЦИЯ: РУКИ, ПРЕДПЛЕЧЬЯ И ИКРЫ
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Подъем штанги на бицепс', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO bicep_curl_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Разгибания на трицепс на верхнем блоке', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO tricep_extension_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Сгибания "Молот" (Хаммер) с гантелями', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO hammer_curl_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Подъемы на носки стоя', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO calf_raise_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (bicep_curl_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'biceps')),
        (tricep_extension_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'triceps')),
        (hammer_curl_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'biceps')),
        (hammer_curl_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'forearms')),
        (calf_raise_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'calves'));

    -- 6. КОР И ПРЕСС
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Планка классическая', 'static', true, def_tenant_id) RETURNING exercise_id INTO plank_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Скручивания на полу', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO crunch_id;
    
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Подъем ног в висе на турнике', 'bodyweight', true, def_tenant_id) RETURNING exercise_id INTO hanging_legs_id;

    INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
        (plank_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'abs')),
        (crunch_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'abs')),
        (hanging_legs_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'abs'));

END $$;