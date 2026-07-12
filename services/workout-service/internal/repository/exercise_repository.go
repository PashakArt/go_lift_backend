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
	//go:embed queries/list_exercises.sql
	listExercisesQuery string

	//go:embed queries/list_exercises_with_user.sql
	listExercisesWithUserIdQuery string
)

type exerciseRepository struct {
	pool *pgxpool.Pool
}

func NewExerciseRepository(pool *pgxpool.Pool) domain.ExerciseRepository {
	return &exerciseRepository{pool: pool}
}

func (r *exerciseRepository) List(ctx context.Context, muscleGroupId uuid.UUID) ([]*domain.Exercise, error) {
	rows, err := r.pool.Query(ctx, listExercisesQuery, muscleGroupId)

	if err != nil {
		return nil, fmt.Errorf("failed to execute list exercises query: %w", err)
	}
	defer rows.Close()

	var exercises []*domain.Exercise
	for rows.Next() {
		var exercise domain.Exercise

		err = rows.Scan(
			&exercise.ExerciseID,
			&exercise.Name,
			&exercise.Type,
			&exercise.IsGlobal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise row: %w", err)
		}
		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercise rows: %w", err)
	}

	return exercises, nil
}

func (r *exerciseRepository) UserFavoriteList(ctx context.Context, muscleGroupId, userId uuid.UUID) ([]*domain.Exercise, error) {
	rows, err := r.pool.Query(ctx, listExercisesWithUserIdQuery, muscleGroupId, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to execute user favorite list exercises query: %w", err)
	}
	defer rows.Close()

	var exercises []*domain.Exercise
	for rows.Next() {
		var exercise domain.Exercise

		err = rows.Scan(
			&exercise.ExerciseID,
			&exercise.Name,
			&exercise.Type,
			&exercise.IsGlobal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise row: %w", err)
		}
		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercise rows: %w", err)
	}

	return exercises, nil
}
