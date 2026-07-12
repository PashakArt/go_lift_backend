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
	List(ctx context.Context, muscleGroupId uuid.UUID) ([]*Exercise, error)
	UserFavoriteList(ctx context.Context, muscleGroupId, userId uuid.UUID) ([]*Exercise, error)
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

type MuscleGroupRepository interface {
	List(ctx context.Context) ([]*MuscleGroup, error)
}
