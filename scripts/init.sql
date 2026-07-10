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
    ('core', 'Кор / Поясница', 80),
    ('glutes', 'Ягодицы', 90),
    ('quads', 'Квадрицепсы (Передняя поверхность бедра)', 100),
    ('hamstrings', 'Бицепс бедра (Задняя поверхность бедра)', 110),
    ('adductors', 'Приводящие мышцы бедра (Внутренняя поверхность)', 120),
    ('calves', 'Икры', 130)
ON CONFLICT (code) DO UPDATE 
SET name = EXCLUDED.name,
    sort_order = EXCLUDED.sort_order;