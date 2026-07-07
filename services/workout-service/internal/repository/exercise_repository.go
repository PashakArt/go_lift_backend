package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/list_exercises.sql
var listExercisesQuery string

type exerciseRepository struct {
	pool *pgxpool.Pool
}

func NewExerciseRepository(pool *pgxpool.Pool) domain.ExerciseRepository {
	return &exerciseRepository{pool: pool}
}

func (r *exerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	return fmt.Errorf("")
}

func (r *exerciseRepository) GetByID(ctx context.Context, tenantID, exerciseID uuid.UUID) (*domain.Exercise, error) {
	return nil, nil
}

func (r *exerciseRepository) List(ctx context.Context, tenantID uuid.UUID, muscleGroupCode string) ([]*domain.Exercise, error) {
	rows, err := r.pool.Query(ctx, listExercisesQuery, tenantID, muscleGroupCode)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list exercises query: %w", err)
	}
	defer rows.Close()

	var exercises []*domain.Exercise
	for rows.Next() {
		var exercise domain.Exercise

		err = rows.Scan(
			&exercise.ExerciseID,
			&exercise.TenantID,
			&exercise.Name,
			&exercise.Type,
			&exercise.IsGlobal,
			&exercise.MuscleGroupCodes,
			&exercise.CreatedAt,
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
