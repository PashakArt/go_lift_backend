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
	Finish(ctx context.Context, userId uuid.UUID) error
}

// WorkoutSetRepository описывает контракт для работы с подходами внутри сессии
type WorkoutSetRepository interface {
	LogWorkoutSet(ctx context.Context, set *WorkoutSet) (*WorkoutSet, error)
}

type MuscleGroupRepository interface {
	List(ctx context.Context) ([]*MuscleGroup, error)
}

type TenantRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*Tenant, error)
}
