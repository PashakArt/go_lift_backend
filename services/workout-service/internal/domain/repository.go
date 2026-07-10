package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository описывает контракт для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByTenantAndTelegramID(ctx context.Context, tenantID uuid.UUID, tgID string) (*User, error)
}

// ExerciseRepository описывает контракт для работы со справочником упражнений
type ExerciseRepository interface {
	Create(ctx context.Context, exercise *Exercise) error
	GetByID(ctx context.Context, tenantID, exerciseID uuid.UUID) (*Exercise, error)
	// List возвращает упражнения арендатора + глобальные упражнения (is_global = true)
	List(ctx context.Context, tenantID uuid.UUID, muscleGroupCode string) ([]*Exercise, error)
}

// WorkoutSessionRepository описывает контракт для управления тренировочными сессиями
type WorkoutSessionRepository interface {
	Create(ctx context.Context, session *WorkoutSession) error
	GetActiveByUserID(ctx context.Context, userId string) (*WorkoutSession, error)
}

// WorkoutSetRepository описывает контракт для работы с подходами внутри сессии
type WorkoutSetRepository interface {
	// CreateMany позволяет сохранить пачку подходов за раз (актуально для круговых тренировок)
	CreateMany(ctx context.Context, sets []*WorkoutSet) error
	// ListBySessionID выгребает все подходы тренировки с сортировкой по sequence_order
	ListBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*WorkoutSet, error)
	// Delete подходов конкретного упражнения в сессии (если юзер решил сбросить прогресс по нему)
	DeleteByExercise(ctx context.Context, sessionID, exerciseID uuid.UUID) error
}
