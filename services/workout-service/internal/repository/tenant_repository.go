package repository

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed queries/get_tenant_by_id.sql
	getTenantById string
)

type tenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) domain.TenantRepository {
	return &tenantRepository{pool}
}

func (r *tenantRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant

	err := r.pool.QueryRow(ctx, getTenantById, id).Scan(
		&tenant.TenantID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute get tenant by id query: %w", err)
	}

	return &tenant, nil
}
