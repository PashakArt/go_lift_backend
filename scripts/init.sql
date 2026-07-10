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
    name VARCHAR(64) NOT NULL     -- 'Грудь', 'Трицепс', 'Бицепс'
);
