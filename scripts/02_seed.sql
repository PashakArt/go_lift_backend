-- ==========================================
-- 1. ДЕФОЛТНЫЙ ТЕНАНТ
-- ==========================================
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
            "logo_url": "https://assets.golift.app/logos/default_logo.png"
        }
    }'::jsonb
)
ON CONFLICT (tenant_id) DO NOTHING;

-- ==========================================
-- 2. СПРАВОЧНИК МЫШЕЧНЫХ ГРУПП
-- ==========================================
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

-- ==========================================
-- 3. БАЗОВЫЕ УПРАЖНЕНИЯ
-- ==========================================
DO $$
DECLARE
    def_tenant_id UUID := '00000000-0000-0000-0000-000000000000';
    
    bench_press_id UUID;
    pushups_id UUID;
    pullups_id UUID;
    barbell_row_id UUID;
    squats_id UUID;
    leg_press_id UUID;
    leg_extension_id UUID;      -- Разгибание ног в тренажере
    adductor_machine_id UUID;
    abductor_machine_id UUID;   -- Разведение ног в тренажере
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

    -- 3. НОГИ, ЯГОДИЦЫ И ПРИВОДЯЩИЕ/ОТВОДЯЩИЕ
    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Приседания со штангой', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO squats_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Жим ногами в тренажере', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO leg_press_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Разгибание ног в тренажере', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO leg_extension_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Сведение ног в тренажере', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO adductor_machine_id;

    INSERT INTO exercises (name, type, is_global, tenant_id) 
    VALUES ('Разведение ног в тренажере', 'dynamic', true, def_tenant_id) RETURNING exercise_id INTO abductor_machine_id;

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
        (leg_extension_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'quads')),
        (adductor_machine_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'adductors')),
        (abductor_machine_id, (SELECT muscle_group_id FROM muscle_groups WHERE code = 'glutes')),
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
