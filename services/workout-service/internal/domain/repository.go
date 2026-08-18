package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository описывает контракт для работы с пользователями
type UserRepository interface {
	Upsert(ctx context.Context, user *User) (bool, error)
	GetByTenantAndTelegramID(ctx context.Context, tenantID uuid.UUID, tgID string) (*User, error)
}

// ExerciseRepository описывает контракт для работы со справочником упражнений
type ExerciseRepository interface {
	List(ctx context.Context, muscleGroupId uuid.UUID) ([]*Exercise, error)
	UserFavoriteList(ctx context.Context, muscleGroupId, userId uuid.UUID) ([]*Exercise, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*Exercise, error)
}

// WorkoutSessionRepository описывает контракт для управления тренировочными сессиями
type WorkoutSessionRepository interface {
	GetWorkoutsForDay(ctx context.Context, userId, date string) ([]WorkoutDaySetRow, error)
	Create(ctx context.Context, session *WorkoutSession) error
	GetActiveByUserID(ctx context.Context, userId string) (*WorkoutSession, error)
	Finish(ctx context.Context, userId uuid.UUID) error
	GetTrainingDays(ctx context.Context, userId string, year, month int) ([]string, error)
	GetUserExportData(ctx context.Context, tgId string) ([]ExportReportRow, error)
}

// WorkoutSetRepository описывает контракт для работы с подходами внутри сессии
type WorkoutSetRepository interface {
	LogWorkoutSet(ctx context.Context, set *WorkoutSet) (*WorkoutSet, error)
	GetCompletedExercises(ctx context.Context, userId, exerciseId string) ([]WorkoutSet, error)
	GetSessionExercises(ctx context.Context, sessionID uuid.UUID) ([]SessionExerciseRow, error)
}

type MuscleGroupRepository interface {
	List(ctx context.Context) ([]*MuscleGroup, error)
}

type TenantRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

type TemplateRepository interface {
	Create(ctx context.Context, template *WorkoutTemplate) (uuid.UUID, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*WorkoutTemplate, error)
	GetByID(ctx context.Context, templateID, userID uuid.UUID) (*WorkoutTemplate, error)
	Delete(ctx context.Context, templateID, userID uuid.UUID) error
	Update(ctx context.Context, template *WorkoutTemplate) error
}
