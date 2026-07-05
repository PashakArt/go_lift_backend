package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository описывает контракт для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, tenantID, userID uuid.UUID) (*User, error)
	GetByTelegramID(ctx context.Context, tenantID uuid.UUID, tgID int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, tenantID, userID uuid.UUID) error
}

// ExerciseRepository описывает контракт для работы со справочником упражнений
type ExerciseRepository interface {
	Create(ctx context.Context, exercise *Exercise) error
	GetByID(ctx context.Context, tenantID, exerciseID uuid.UUID) (*Exercise, error)
	// List возвращает упражнения арендатора + глобальные упражнения (is_global = true)
	List(ctx context.Context, tenantID uuid.UUID, muscleGroup string) ([]*Exercise, error)
	Update(ctx context.Context, exercise *Exercise) error
	Delete(ctx context.Context, tenantID, exerciseID uuid.UUID) error
}

// WorkoutSessionRepository описывает контракт для управления тренировочными сессиями
type WorkoutSessionRepository interface {
	Create(ctx context.Context, session *WorkoutSession) error
	GetByID(ctx context.Context, tenantID, sessionID uuid.UUID) (*WorkoutSession, error)
	// GetActiveSession возвращает текущую незавершенную тренировку пользователя (где ended_at IS NULL)
	GetActiveSession(ctx context.Context, tenantID, userID uuid.UUID) (*WorkoutSession, error)
	// ListByUserID возвращает историю тренировок пользователя с пагинацией (limit, offset)
	ListByUserID(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]*WorkoutSession, error)
	Update(ctx context.Context, session *WorkoutSession) error
	Delete(ctx context.Context, tenantID, sessionID uuid.UUID) error
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
