package repository

import (
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User        domain.UserRepository
	Exercise    domain.ExerciseRepository
	Session     domain.WorkoutSessionRepository
	MuscleGroup domain.MuscleGroupRepository
	Tenant      domain.TenantRepository
	WorkoutSet  domain.WorkoutSetRepository
	Template    domain.TemplateRepository
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:        NewUserRepository(pool),
		Exercise:    NewExerciseRepository(pool),
		Session:     NewWorkoutSessionRepository(pool),
		MuscleGroup: NewMuscleGroupRepository(pool),
		Tenant:      NewTenantRepository(pool),
		WorkoutSet:  NewWWorkoutSetRepository(pool),
		Template:    NewTemplateRepository(pool),
	}
}
