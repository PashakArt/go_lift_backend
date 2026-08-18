package repository

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

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

	//go:embed queries/get_training_days.sql
	getTrainingDaysQuery string

	//go:embed queries/get_workouts_for_day.sql
	getWorkoutsForDayQuery string

	//go:embed queries/get_user_export_data.sql
	getUserExportData string
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
		&session.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute get active session by user_id query: %w", err)
	}

	return &session, nil
}

func (s *workoutSessionRepository) GetTrainingDays(ctx context.Context, userId string, year, month int) ([]string, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	rows, err := s.pool.Query(ctx, getTrainingDaysQuery, userId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get training days query: %w", err)
	}
	defer rows.Close()

	days := make([]string, 0)

	for rows.Next() {
		var day string
		if err := rows.Scan(&day); err != nil {
			return nil, fmt.Errorf("failed to scan training day: %w", err)
		}
		days = append(days, day)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return days, nil
}

func (r *workoutSessionRepository) GetWorkoutsForDay(ctx context.Context, userId, dateStr string) ([]domain.WorkoutDaySetRow, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	rows, err := r.pool.Query(ctx, getWorkoutsForDayQuery, userId, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get_workouts_for_day query: %w", err)
	}
	defer rows.Close()

	var result []domain.WorkoutDaySetRow
	for rows.Next() {
		var row domain.WorkoutDaySetRow
		err := rows.Scan(
			&row.WorkoutSession.SessionID,
			&row.StartedAt,
			&row.EndedAt,
			&row.Type,
			&row.SetID,
			&row.ExerciseID,
			&row.ExerciseName,
			&row.ExerciseType,
			&row.SetNumber,
			&row.Weight,
			&row.Reps,
			&row.DurationSeconds,
			&row.DistanceMeters,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan WorkoutDaySetRow: %w", err)
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func (r *workoutSessionRepository) GetUserExportData(ctx context.Context, tgId string) ([]domain.ExportReportRow, error) {
	rows, err := r.pool.Query(ctx, getUserExportData, tgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ExportReportRow
	for rows.Next() {
		var row domain.ExportReportRow
		if err := rows.Scan(
			&row.SessionId,
			&row.StartedAt,
			&row.EndedAt,
			&row.SessionType,
			&row.ExerciseName,
			&row.SetNumber,
			&row.Weight,
			&row.Reps,
			&row.DurationSec,
			&row.DistanceM,
		); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}
