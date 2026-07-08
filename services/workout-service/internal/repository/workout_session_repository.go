package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/start_session.sql
var startSessionQuery string

type workoutSessionRepository struct {
	pool *pgxpool.Pool
}

func NewWorkoutSessionRepository(pool *pgxpool.Pool) domain.WorkoutSessionRepository {
	return &workoutSessionRepository{pool: pool}
}

func (r *workoutSessionRepository) Create(ctx context.Context, session *domain.WorkoutSession) error {
	var templateStr string
	if session.TemplateID != nil {
		templateStr = session.TemplateID.String()
	}

	err := r.pool.QueryRow(
		ctx,
		startSessionQuery,
		session.TenantID,
		session.UserID,
		templateStr,
		session.Type,
	).Scan(&session.SessionID, &session.StartedAt)

	if err != nil {
		return fmt.Errorf("failed to execute start session query: %w", err)
	}

	return nil
}
