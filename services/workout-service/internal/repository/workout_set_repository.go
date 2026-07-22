package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed queries/log_set_query.sql
	logWorkoutSetQuery string

	//go:embed queries/get_completed_exercises.sql
	getCompletedExercisesQuery string
)

type workoutSetRepository struct {
	pool *pgxpool.Pool
}

func NewWWorkoutSetRepository(pool *pgxpool.Pool) domain.WorkoutSetRepository {
	return &workoutSetRepository{pool: pool}
}

func (r *workoutSetRepository) LogWorkoutSet(ctx context.Context, set *domain.WorkoutSet) (*domain.WorkoutSet, error) {
	if set.SetID == uuid.Nil {
		set.SetID = uuid.New()
	}

	res := &domain.WorkoutSet{}

	err := r.pool.QueryRow(
		ctx,
		logWorkoutSetQuery,
		set.SetID,
		set.SessionID,
		set.ExerciseID,
		set.TenantID,
		set.SetNumber,
		set.Weight,
		set.Reps,
		set.DurationSeconds,
		set.DistanceMeters,
	).Scan(
		&res.SetID,
		&res.SessionID,
		&res.ExerciseID,
		&res.TenantID,
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

func (r *workoutSetRepository) GetCompletedExercises(ctx context.Context, userId, exerciseId string) ([]domain.WorkoutSet, error) {
	rows, err := r.pool.Query(ctx, getCompletedExercisesQuery, userId, exerciseId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetCompletedExercises query: %w", err)
	}
	defer rows.Close()

	var exercises []domain.WorkoutSet
	for rows.Next() {
		var exercise domain.WorkoutSet

		err = rows.Scan(
			&exercise.SetNumber,
			&exercise.Weight,
			&exercise.Reps,
			&exercise.DurationSeconds,
			&exercise.DistanceMeters,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise row: %w", err)
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}
