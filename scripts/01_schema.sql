CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- 1. ТЕНАНТЫ И ПОЛЬЗОВАТЕЛИ
-- ==========================================

CREATE TABLE IF NOT EXISTS tenants (
    tenant_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(100) NOT NULL UNIQUE, -- Добавлено прямо в определение таблицы
    branding_json JSONB,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    user_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    telegram_id VARCHAR(64) NOT NULL,
    phone       VARCHAR(32),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_tenant_telegram UNIQUE (tenant_id, telegram_id)
);

-- ==========================================
-- 2. СПРАВОЧНИКИ УПРАЖНЕНИЙ И МЫШЦ
-- ==========================================

CREATE TABLE IF NOT EXISTS muscle_groups (
    muscle_group_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(32) NOT NULL UNIQUE,
    name            VARCHAR(64) NOT NULL,
    sort_order      INT NOT NULL DEFAULT 0
);

-- Безопасное создание ENUM типом (не будет падать, если тип уже существует)
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exercise_type') THEN
        CREATE TYPE exercise_type AS ENUM ('dynamic', 'static', 'bodyweight', 'cardio');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS exercises (
    exercise_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    type        exercise_type NOT NULL,
    is_global   BOOLEAN DEFAULT false,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS exercise_muscle_groups (
    exercise_id     UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
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
-- 3. ТРЕНИРОВОЧНЫЕ СЕССИИ И ТЕМПЛЕЙТЫ
-- ==========================================

CREATE TABLE IF NOT EXISTS workout_templates (
    template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    creator_id  UUID REFERENCES users(user_id) ON DELETE SET NULL,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'session_type') THEN
        CREATE TYPE session_type AS ENUM ('classic', 'circuit');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS workout_sessions (
    session_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id     UUID REFERENCES users(user_id) ON DELETE CASCADE,
    template_id UUID REFERENCES workout_templates(template_id) ON DELETE SET NULL,
    type        session_type NOT NULL DEFAULT 'classic',
    started_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ended_at    TIMESTAMP WITH TIME ZONE,
    is_active   BOOLEAN DEFAULT true
);

CREATE UNIQUE INDEX idx_unique_active_session 
ON workout_sessions (user_id) 
WHERE is_active = true;

CREATE TABLE IF NOT EXISTS workout_sets (
    set_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id       UUID NOT NULL REFERENCES workout_sessions(session_id) ON DELETE CASCADE,
    exercise_id      UUID NOT NULL REFERENCES exercises(exercise_id),
    tenant_id        UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    
    set_number       INT NOT NULL,
    round_number     INT NOT NULL DEFAULT 1,

    weight           NUMERIC(6, 2),
    reps             INT,
    duration_seconds INT,
    distance_meters  INT,
    
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- <-- Исправлено (добавлена запятая)

    CONSTRAINT uq_workout_sets_session_exercise_set_round 
        UNIQUE (session_id, exercise_id, set_number, round_number)
);

CREATE INDEX IF NOT EXISTS idx_workout_sets_session ON workout_sets(session_id);
CREATE INDEX IF NOT EXISTS idx_workout_sets_exercise_user ON workout_sets(tenant_id, exercise_id);