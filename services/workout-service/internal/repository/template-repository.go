package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/create_template.sql
var createTemplateQuery string

type templateRepository struct {
	pool *pgxpool.Pool
}

func NewTemplateRepository(pool *pgxpool.Pool) domain.TemplateRepository {
	return &templateRepository{pool: pool}
}

func (r *templateRepository) Create(ctx context.Context, template *domain.WorkoutTemplate) (uuid.UUID, error) {
	itemsJSON, err := json.Marshal(template.Items)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal template items to json: %w", err)
	}

	var templateID uuid.UUID
	err = r.pool.QueryRow(
		ctx,
		createTemplateQuery,
		template.UserID,
		template.Name,
		itemsJSON,
	).Scan(&templateID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert workout template: %w", err)
	}

	return templateID, nil
}
