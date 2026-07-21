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
	//go:embed queries/start_session.sql
	startSessionQuery string

	//go:embed queries/finish_session.sql
	finishSessionQuery string

	//go:embed queries/get_active_session_by_user_id.sql
	getActiveSessionByUserId string
)

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

func (r *workoutSessionRepository) Finish(ctx context.Context, userId uuid.UUID) error {
	_, err := r.pool.Exec(ctx, finishSessionQuery, userId)
	if err != nil {
		return fmt.Errorf("failed to execute finish session query: %w", err)
	}

	return nil
}

func (r *workoutSessionRepository) GetActiveByUserID(ctx context.Context, userId string) (*domain.WorkoutSession, error) {
	var session domain.WorkoutSession

	err := r.pool.QueryRow(ctx, getActiveSessionByUserId, userId).Scan(
		&session.SessionID,
		&session.TenantID,
		&session.UserID,
		&session.TemplateID,
		&session.Type,
		&session.StartedAt,
		&session.EndedAt,
		&session.IsACtive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute get active session by user_id query: %w", err)
	}

	return &session, nil
}
