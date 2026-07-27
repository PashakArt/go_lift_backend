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

var (
	//go:embed queries/create_template.sql
	createTemplateQuery string

	//go:embed queries/get_templates.sql
	getTemplatesQuery string

	//go:embed queries/get_template_by_id.sql
	getTemplateByIDQuery string

	//go:embed queries/delete_template.sql
	deleteTemplateQuery string

	//go:embed queries/update_template.sql
	updateTemplateQuery string
)

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

func (r *templateRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.WorkoutTemplate, error) {
	rows, err := r.pool.Query(ctx, getTemplatesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*domain.WorkoutTemplate

	for rows.Next() {
		var (
			t         domain.WorkoutTemplate
			itemsJSON []byte
		)

		err := rows.Scan(
			&t.TemplateID,
			&t.UserID,
			&t.Name,
			&itemsJSON,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template row: %w", err)
		}

		if len(itemsJSON) > 0 {
			if err := json.Unmarshal(itemsJSON, &t.Items); err != nil {
				return nil, fmt.Errorf("failed to unmarshal template items: %w", err)
			}
		}

		templates = append(templates, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return templates, nil
}

func (r *templateRepository) GetByID(ctx context.Context, templateID, userID uuid.UUID) (*domain.WorkoutTemplate, error) {
	var (
		t         domain.WorkoutTemplate
		itemsJSON []byte
	)

	err := r.pool.QueryRow(ctx, getTemplateByIDQuery, templateID, userID).Scan(
		&t.TemplateID,
		&t.UserID,
		&t.Name,
		&itemsJSON,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan template by id: %w", err)
	}

	if len(itemsJSON) > 0 {
		if err := json.Unmarshal(itemsJSON, &t.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template items: %w", err)
		}
	}

	return &t, nil
}

func (r *templateRepository) Delete(ctx context.Context, templateID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, deleteTemplateQuery, templateID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete workout template: %w", err)
	}

	return nil
}

func (r *templateRepository) Update(ctx context.Context, template *domain.WorkoutTemplate) error {
	itemsJSON, err := json.Marshal(template.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal template items to json: %w", err)
	}

	cmdTag, err := r.pool.Exec(
		ctx,
		updateTemplateQuery,
		template.Name,
		itemsJSON,
		template.TemplateID,
		template.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update workout template: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("workout template not found or access denied")
	}

	return nil
}
