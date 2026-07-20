package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/log_set_query.sql
var logWorkoutSetQuery string

type WorkoutSetRepository struct {
	pool *pgxpool.Pool
}

func NewWWorkoutSetRepository(pool *pgxpool.Pool) domain.WorkoutSetRepository {
	return &workoutSessionRepository{pool: pool}
}

func (r *workoutSessionRepository) LogWorkoutSet(ctx context.Context, set *domain.WorkoutSet) (*domain.WorkoutSet, error) {
	var setIDParam interface{}
	if set.SetID != uuid.Nil {
		setIDParam = set.SetID.String()
	} else {
		setIDParam = nil
	}

	res := &domain.WorkoutSet{}

	err := r.pool.QueryRow(
		ctx,
		logWorkoutSetQuery,
		setIDParam,
		set.SessionID,
		set.ExerciseID,
		set.SetNumber,
		set.Weight,
		set.Reps,
		set.DurationSeconds,
		set.DistanceMeters,
	).Scan(
		&res.SetID,
		&res.SessionID,
		&res.ExerciseID,
		&res.SetNumber,
		&res.Weight,
		&res.Reps,
		&res.DurationSeconds,
		&res.DistanceMeters,
		&res.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute LogWorkoutSet query: %w", err)
	}

	return res, nil
}
