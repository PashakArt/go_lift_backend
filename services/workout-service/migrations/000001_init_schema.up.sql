-- =========================================================================
-- 1. ТАБЛИЦА АРЕНДАТОРОВ (TENANTS)
-- =========================================================================
CREATE TABLE tenants (
    tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_token_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    branding_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE tenants IS 'Таблица арендаторов (фитнес-студий или розничных контуров платформы) для реализации Multi-tenant архитектуры';
COMMENT ON COLUMN tenants.tenant_id IS 'Уникальный UUID идентификатор арендатора';
COMMENT ON COLUMN tenants.bot_token_hash IS 'Криптографический хэш токена Telegram-бота для безопасной валидации вебхуков';
COMMENT ON COLUMN tenants.name IS 'Коммерческое название фитнес-студии или обозначение базового контура GoLift Retail';
COMMENT ON COLUMN tenants.branding_json IS 'Конфигурационный JSON для фронтенда Telegram Mini App (фирменные цвета, ссылки на логотипы, стили)';
COMMENT ON COLUMN tenants.created_at IS 'Временная метка создания записи арендатора';


-- =========================================================================
-- 2. ТАБЛИЦА ПОЛЬЗОВАТЕЛЕЙ (USERS)
-- =========================================================================
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE SET NULL,
    telegram_id BIGINT NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE (tenant_id, telegram_id)
);

COMMENT ON TABLE users IS 'Пользователи системы (атлеты и тренеры), привязанные к конкретным арендаторам';
COMMENT ON COLUMN users.user_id IS 'Уникальный UUID идентификатор пользователя';
COMMENT ON COLUMN users.tenant_id IS 'Ссылка на арендатора. NULL означает независимого розничного пользователя';
COMMENT ON COLUMN users.telegram_id IS 'Оригинальный числовой ID пользователя, полученный из Telegram API';
COMMENT ON COLUMN users.phone IS 'Номер телефона пользователя (опционально, для B2B-синхронизации с CRM залов)';
COMMENT ON COLUMN users.created_at IS 'Временная метка регистрации пользователя в системе';


-- =========================================================================
-- 3. ТАБЛИЦА УПРАЖНЕНИЙ (EXERCISES)
-- =========================================================================
CREATE TYPE exercise_type AS ENUM ('dynamic', 'static', 'bodyweight', 'cardio');
COMMENT ON TYPE exercise_type IS 'Полиморфные типы нагрузок: dynamic (с весом), static (на время), bodyweight (свой вес), cardio (дистанция+время)';

CREATE TABLE exercises (
    exercise_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type exercise_type NOT NULL,
    muscle_group VARCHAR(50) NOT NULL,
    is_global BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE exercises IS 'Справочник упражнений. Содержит как глобальные записи, так и кастомные упражнения конкретных студий';
COMMENT ON COLUMN exercises.exercise_id IS 'Уникальный UUID идентификатор упражнения';
COMMENT ON COLUMN exercises.tenant_id IS 'Ссылка на создателя упражнения. Для глобальных упражнений может указывать на системный tenant';
COMMENT ON COLUMN exercises.name IS 'Название упражнения (например, Жим штанги лежа)';
COMMENT ON COLUMN exercises.type IS 'Тип физической нагрузки, определяющий набор метрик в таблице подходов';
COMMENT ON COLUMN exercises.muscle_group IS 'Целевая мышечная группа (грудь, спина, ноги, плечи, руки, кор)';
COMMENT ON COLUMN exercises.is_global IS 'Флаг, определяющий, доступно ли упражнение всем пользователям платформы вне зависимости от их tenant_id';
COMMENT ON COLUMN exercises.created_at IS 'Временная метка добавления упражнения в базу';


-- =========================================================================
-- 4. ТАБЛИЦА ШАБЛОНОВ ТРЕНИРОВОК (WORKOUT TEMPLATES)
-- =========================================================================
CREATE TABLE workout_templates (
    template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    creator_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE workout_templates IS 'Шаблоны готовых программ тренировок, созданные пользователями или тренерами фитнес-студий';
COMMENT ON COLUMN workout_templates.template_id IS 'Уникальный UUID идентификатор шаблона';
COMMENT ON COLUMN workout_templates.tenant_id IS 'Ссылка на арендатора, внутри которого доступен данный шаблон';
COMMENT ON COLUMN workout_templates.creator_id IS 'Ссылка на пользователя (атлета или тренера), который спроектировал этот шаблон';
COMMENT ON COLUMN workout_templates.name IS 'Название программы (например, День ног, На массу А)';
COMMENT ON COLUMN workout_templates.description IS 'Текстовое описание или методические указания к тренировочному комплексу';
COMMENT ON COLUMN workout_templates.created_at IS 'Временная метка создания шаблона';


-- Связующая таблица для структуры упражнений в шаблоне
CREATE TABLE template_exercises (
    template_exercise_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES workout_templates(template_id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    sequence_order INT NOT NULL,
    UNIQUE (template_id, sequence_order)
);

COMMENT ON TABLE template_exercises IS 'Связующая таблица для фиксации состава упражнений внутри шаблона с сохранением строгого порядка выполнения';
COMMENT ON COLUMN template_exercises.template_exercise_id IS 'Уникальный UUID идентификатор связи';
COMMENT ON COLUMN template_exercises.template_id IS 'Ссылка на целевой шаблон тренировки';
COMMENT ON COLUMN template_exercises.exercise_id IS 'Ссылка на добавляемое упражнение из справочника';
COMMENT ON COLUMN template_exercises.sequence_order IS 'Порядковый номер выполнения упражнения в рамках данного шаблона (1, 2, 3...)';


-- =========================================================================
-- 5. ТАБЛИЦА СЕССИЙ ТРЕНИРОВОК (WORKOUT SESSIONS)
-- =========================================================================
CREATE TYPE session_type AS ENUM ('classic', 'circuit');
COMMENT ON TYPE session_type IS 'Режим выполнения тренировки: classic (последовательный), circuit (круговой кроссфит-круг)';

CREATE TABLE workout_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    template_id UUID REFERENCES workout_templates(template_id) ON DELETE SET NULL,
    type session_type NOT NULL DEFAULT 'classic',
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    notes TEXT
);

COMMENT ON TABLE workout_sessions IS 'Исторический лог запущенных и выполненных тренировочных сессий пользователей';
COMMENT ON COLUMN workout_sessions.session_id IS 'Уникальный UUID идентификатор конкретной тренировочной сессии';
COMMENT ON COLUMN workout_sessions.tenant_id IS 'Ссылка на арендатора, под эгидой которого проходит тренировка';
COMMENT ON COLUMN workout_sessions.user_id IS 'Ссылка на пользователя, выполняющего тренировку';
COMMENT ON COLUMN workout_sessions.template_id IS 'Ссылка на исходный шаблон, если тренировка была запущена по программе. Может быть NULL';
COMMENT ON COLUMN workout_sessions.type IS 'Режим тренировки, определяющий логику переключения экранов в Mini App и группировку подходов';
COMMENT ON COLUMN workout_sessions.started_at IS 'Время фактического старта тренировочной сессии';
COMMENT ON COLUMN workout_sessions.ended_at IS 'Время завершения тренировки. Заполняется при нажатии кнопки Закончить';
COMMENT ON COLUMN workout_sessions.notes IS 'Финальный текстовый комментарий пользователя или тренера по итогам всей сессии';


-- =========================================================================
-- 6. ТАБЛИЦА ПОДХОДОВ / РЕЗУЛЬТАТОВ (WORKOUT SETS)
-- =========================================================================
CREATE TABLE workout_sets (
    set_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES workout_sessions(session_id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    round_number INT NOT NULL DEFAULT 1,
    sequence_order INT NOT NULL,
    weight NUMERIC(6, 2),
    reps INT,
    duration_seconds INT,
    distance_meters INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE workout_sets IS 'Единая полиморфная таблица для логирования каждого выполненного подхода/раунда во всех типах физических нагрузок';
COMMENT ON COLUMN workout_sets.set_id IS 'Уникальный UUID идентификатор конкретного подхода';
COMMENT ON COLUMN workout_sets.session_id IS 'Ссылка на сессию тренировки, в рамках которой выполнен подход';
COMMENT ON COLUMN workout_sets.exercise_id IS 'Ссылка на выполняемое упражнение';
COMMENT ON COLUMN workout_sets.round_number IS 'Номер круга. Для классического режима равен 1. Для круговых тренировок увеличивается с каждым новым циклом упражнений';
COMMENT ON COLUMN workout_sets.sequence_order IS 'Сквозной порядковый номер записи внутри тренировки или раунда для обеспечения корректной хронологии на графиках';
COMMENT ON COLUMN workout_sets.weight IS 'Метрика динамического типа: вес отягощения в килограммах (например, 82.50)';
COMMENT ON COLUMN workout_sets.reps IS 'Метрика динамического и собственного веса: количество успешно выполненных повторений';
COMMENT ON COLUMN workout_sets.duration_seconds IS 'Метрика статического и кардио типа: время удержания позиции или выполнения аэробной работы в секундах';
COMMENT ON COLUMN workout_sets.distance_meters IS 'Метрика кардио типа: преодоленное расстояние в метрах';
COMMENT ON COLUMN workout_sets.created_at IS 'Точное время фиксации выполнения конкретного подхода (используется бэкендом для расчета времени отдыха)';

-- ИНДЕКСЫ ДЛЯ УСКОРЕНИЯ ВЫБЕРОК
CREATE INDEX idx_workout_sets_exercise_session ON workout_sets(exercise_id, session_id);
CREATE INDEX idx_sessions_user_tenant ON workout_sessions(user_id, tenant_id);