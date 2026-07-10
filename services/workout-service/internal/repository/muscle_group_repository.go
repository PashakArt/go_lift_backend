package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/list_muscle_groups.sql
var listMuscleGroupsQuery string

type muscleGroupRepository struct {
	pool *pgxpool.Pool
}

func NewMuscleGroupRepository(pool *pgxpool.Pool) domain.MuscleGroupRepository {
	return &muscleGroupRepository{pool}
}

func (r *muscleGroupRepository) List(ctx context.Context) ([]*domain.MuscleGroup, error) {
	rows, err := r.pool.Query(ctx, listMuscleGroupsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list muscle group query: %w", err)
	}
	defer rows.Close()

	var muscleGroups []*domain.MuscleGroup
	for rows.Next() {
		var mg domain.MuscleGroup

		err = rows.Scan(
			&mg.MuscleGroupId,
			&mg.Code,
			&mg.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise row: %w", err)
		}
		muscleGroups = append(muscleGroups, &mg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercise rows: %w", err)
	}

	return muscleGroups, nil
}
