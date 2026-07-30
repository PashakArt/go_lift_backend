DO $$
DECLARE
    def_tenant_id UUID := '00000000-0000-0000-0000-000000000000';

    -- Переменные ID новых упражнений
    matrix_bench_press_id UUID;
    seated_chest_press_id UUID;
    pec_deck_fly_id UUID;
    seated_row_neutral_id UUID;
    seated_row_unilateral_id UUID;
    romanian_deadlift_id UUID;
    seated_leg_curl_id UUID;

    chest_mg_id UUID;
    triceps_mg_id UUID;
    shoulders_mg_id UUID;
    back_mg_id UUID;
    biceps_mg_id UUID;
    forearms_mg_id UUID;
    hamstrings_mg_id UUID;
    glutes_mg_id UUID;
BEGIN
    SELECT muscle_group_id INTO chest_mg_id FROM muscle_groups WHERE code = 'chest';
    SELECT muscle_group_id INTO triceps_mg_id FROM muscle_groups WHERE code = 'triceps';
    SELECT muscle_group_id INTO shoulders_mg_id FROM muscle_groups WHERE code = 'shoulders';
    SELECT muscle_group_id INTO back_mg_id FROM muscle_groups WHERE code = 'back';
    SELECT muscle_group_id INTO biceps_mg_id FROM muscle_groups WHERE code = 'biceps';
    SELECT muscle_group_id INTO forearms_mg_id FROM muscle_groups WHERE code = 'forearms';
    SELECT muscle_group_id INTO hamstrings_mg_id FROM muscle_groups WHERE code = 'hamstrings';
    SELECT muscle_group_id INTO glutes_mg_id FROM muscle_groups WHERE code = 'glutes';

    -- ==========================================
    -- 1. Жим лежа в тренажере Matrix
    -- ==========================================
    SELECT exercise_id INTO matrix_bench_press_id FROM exercises WHERE name = 'Жим лежа в тренажере Matrix' AND tenant_id = def_tenant_id;
    IF matrix_bench_press_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Жим лежа в тренажере Matrix', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO matrix_bench_press_id;
        
        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (matrix_bench_press_id, chest_mg_id),
            (matrix_bench_press_id, triceps_mg_id),
            (matrix_bench_press_id, shoulders_mg_id);
    END IF;

    -- ==========================================
    -- 2. Жим от груди сидя в тренажере
    -- ==========================================
    SELECT exercise_id INTO seated_chest_press_id FROM exercises WHERE name = 'Жим от груди сидя в тренажере' AND tenant_id = def_tenant_id;
    IF seated_chest_press_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Жим от груди сидя в тренажере', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO seated_chest_press_id;

        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (seated_chest_press_id, chest_mg_id),
            (seated_chest_press_id, triceps_mg_id),
            (seated_chest_press_id, shoulders_mg_id);
    END IF;

    -- ==========================================
    -- 3. Сведение рук в тренажере (Бабочка / Pec-Deck)
    -- ==========================================
    SELECT exercise_id INTO pec_deck_fly_id FROM exercises WHERE name = 'Сведение рук в тренажере Matrix (Бабочка)' AND tenant_id = def_tenant_id;
    IF pec_deck_fly_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Сведение рук в тренажере Matrix (Бабочка)', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO pec_deck_fly_id;

        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (pec_deck_fly_id, chest_mg_id);
    END IF;

    -- ==========================================
    -- 4. Горизонтальная тяга в тренажере нейтральным хватом
    -- ==========================================
    SELECT exercise_id INTO seated_row_neutral_id FROM exercises WHERE name = 'Горизонтальная тяга в тренажере нейтральным хватом' AND tenant_id = def_tenant_id;
    IF seated_row_neutral_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Горизонтальная тяга в тренажере нейтральным хватом', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO seated_row_neutral_id;

        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (seated_row_neutral_id, back_mg_id),
            (seated_row_neutral_id, biceps_mg_id);
    END IF;

    -- ==========================================
    -- 5. Горизонтальная тяга в тренажере одной рукой (поочередно)
    -- ==========================================
    SELECT exercise_id INTO seated_row_unilateral_id FROM exercises WHERE name = 'Горизонтальная тяга в тренажере одной рукой' AND tenant_id = def_tenant_id;
    IF seated_row_unilateral_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Горизонтальная тяга в тренажере одной рукой', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO seated_row_unilateral_id;

        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (seated_row_unilateral_id, back_mg_id),
            (seated_row_unilateral_id, biceps_mg_id),
            (seated_row_unilateral_id, forearms_mg_id);
    END IF;

    -- ==========================================
    -- 6. Румынская тяга со штангой
    -- ==========================================
    SELECT exercise_id INTO romanian_deadlift_id FROM exercises WHERE name = 'Румынская тяга со штангой' AND tenant_id = def_tenant_id;
    IF romanian_deadlift_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Румынская тяга со штангой', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO romanian_deadlift_id;

        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (romanian_deadlift_id, hamstrings_mg_id),
            (romanian_deadlift_id, glutes_mg_id),
            (romanian_deadlift_id, back_mg_id);
    END IF;

    -- ==========================================
    -- 7. Сгибание ног сидя в тренажере
    -- ==========================================
    SELECT exercise_id INTO seated_leg_curl_id FROM exercises WHERE name = 'Сгибание ног сидя в тренажере' AND tenant_id = def_tenant_id;
    IF seated_leg_curl_id IS NULL THEN
        INSERT INTO exercises (name, type, is_global, tenant_id)
        VALUES ('Сгибание ног сидя в тренажере', 'dynamic', true, def_tenant_id)
        RETURNING exercise_id INTO seated_leg_curl_id;

        INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id) VALUES 
            (seated_leg_curl_id, hamstrings_mg_id);
    END IF;

END $$;