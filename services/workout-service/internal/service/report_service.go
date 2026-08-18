package service

import (
	"context"
	"fmt"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/exporter"
	"github.com/google/uuid"
)

type ReportService interface {
	ExportWorkoutReport(ctx context.Context, tgId string) ([]byte, error)
}

type reportService struct {
	sessionRepo domain.WorkoutSessionRepository
	exporter    *exporter.ExcelExporter
}

func NewReportService(
	sessionRepo domain.WorkoutSessionRepository,
	exp *exporter.ExcelExporter,
) ReportService {
	return &reportService{
		sessionRepo: sessionRepo,
		exporter:    exp,
	}
}

func (s *reportService) ExportWorkoutReport(ctx context.Context, tgId string) ([]byte, error) {
	reportRows, err := s.sessionRepo.GetUserExportData(ctx, tgId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user export data: %w", err)
	}

	sessionMap := make(map[uuid.UUID]*exporter.WorkoutReportSession)
	var sessionOrder []uuid.UUID

	for _, r := range reportRows {
		sess, exists := sessionMap[r.SessionId]
		if !exists {
			sess = &exporter.WorkoutReportSession{
				StartedAt:   r.StartedAt,
				EndedAt:     r.EndedAt,
				SessionType: r.SessionType,
				Exercises:   []exporter.WorkoutReportSet{},
			}
			sessionMap[r.SessionId] = sess
			sessionOrder = append(sessionOrder, r.SessionId)
		}

		var reps *int32
		if r.Reps != nil {
			v := int32(*r.Reps)
			reps = &v
		}

		var dur *int32
		if r.DurationSec != nil {
			v := int32(*r.DurationSec)
			dur = &v
		}

		var dist *int32
		if r.DistanceM != nil {
			v := int32(*r.DistanceM)
			dist = &v
		}

		var weight *float64
		if r.Weight != nil {
			v := float64(*r.Weight)
			weight = &v
		}

		sess.Exercises = append(sess.Exercises, exporter.WorkoutReportSet{
			SetNumber:    r.SetNumber,
			ExerciseName: r.ExerciseName,
			Weight:       weight,
			Reps:         reps,
			DurationSec:  dur,
			DistanceM:    dist,
		})
	}

	sessions := make([]exporter.WorkoutReportSession, 0, len(sessionOrder))
	for _, id := range sessionOrder {
		sessions = append(sessions, *sessionMap[id])
	}

	return s.exporter.GenerateWorkoutReport(sessions)
}
