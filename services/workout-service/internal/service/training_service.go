package service

import (
	"context"
	"fmt"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/types"
	"github.com/google/uuid"
)

type TrainingService interface {
	StartSession(ctx context.Context, tenantIDStr, userIDStr, sessionType, templateIDStr string) (*domain.WorkoutSession, error)
	LogWorkoutSet(ctx context.Context, req *workoutv1.LogSetRequest) (*domain.WorkoutSet, error)
	FinishSession(ctx context.Context, userId string) error
	GetCompletedExercises(ctx context.Context, userId, exerciseId string) ([]domain.WorkoutSet, error)
	GetTrainingDays(ctx context.Context, userId string, year, month int) ([]string, error)
	GetWorkoutsForDay(ctx context.Context, userId, date string) (*types.WorkoutForDay, error)
	GetSessionExercises(ctx context.Context, sessionID string) ([]types.WorkoutExercise, error)
}

type trainingService struct {
	sessionRepo    domain.WorkoutSessionRepository
	workoutSetRepo domain.WorkoutSetRepository
}

func NewTrainingService(s domain.WorkoutSessionRepository, w domain.WorkoutSetRepository) TrainingService {
	return &trainingService{
		sessionRepo:    s,
		workoutSetRepo: w,
	}
}

func (s *trainingService) StartSession(ctx context.Context, tenantIDStr, userIDStr, sessionType, templateIDStr string) (*domain.WorkoutSession, error) {
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id format: %w", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	var templateID *uuid.UUID
	if templateIDStr != "" {
		parsedTemplateID, err := uuid.Parse(templateIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid template id format: %w", err)
		}
		templateID = &parsedTemplateID
	}

	session := &domain.WorkoutSession{
		TenantID:   tenantID,
		UserID:     userID,
		TemplateID: templateID,
		Type:       domain.SessionType(sessionType),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create workout session in repo: %w", err)
	}

	return session, nil
}

func (s *trainingService) FinishSession(ctx context.Context, userIDStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	if err := s.sessionRepo.Finish(ctx, userID); err != nil {
		return fmt.Errorf("failed to finish workout session in repo: %w", err)
	}

	return nil
}

func (s *trainingService) LogWorkoutSet(
	ctx context.Context,
	protoReq *workoutv1.LogSetRequest,
) (*domain.WorkoutSet, error) {
	sessionID, err := uuid.Parse(protoReq.GetSessionId())
	if err != nil {
		return nil, fmt.Errorf("invalid session id: %w", err)
	}

	exerciseID, err := uuid.Parse(protoReq.GetExerciseId())
	if err != nil {
		return nil, fmt.Errorf("invalid exercise id: %w", err)
	}

	tenantID, err := uuid.Parse(protoReq.GetTenantId())
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}

	var setID uuid.UUID
	if protoReq.SetId != nil && *protoReq.SetId != "" {
		parsedSetID, err := uuid.Parse(*protoReq.SetId)
		if err != nil {
			return nil, fmt.Errorf("invalid set id format: %w", err)
		}
		setID = parsedSetID
	}

	var weight *float64
	if protoReq.Weight != nil {
		w := float64(*protoReq.Weight)
		weight = &w
	}

	var reps *int
	if protoReq.Reps != nil {
		r := int(*protoReq.Reps)
		reps = &r
	}

	var durationSec *int
	if protoReq.DurationSec != nil {
		d := int(*protoReq.DurationSec)
		durationSec = &d
	}

	var distanceMet *int
	if protoReq.DistanceMet != nil {
		dm := int(*protoReq.DistanceMet)
		distanceMet = &dm
	}

	set := &domain.WorkoutSet{
		SetID:           setID,
		SessionID:       sessionID,
		ExerciseID:      exerciseID,
		TenantID:        tenantID,
		Weight:          weight,
		Reps:            reps,
		DurationSeconds: durationSec,
		DistanceMeters:  distanceMet,
	}

	res, err := s.workoutSetRepo.LogWorkoutSet(ctx, set)
	if err != nil {
		return nil, fmt.Errorf("failed to log workout set: %w", err)
	}

	return res, nil
}

func (s *trainingService) GetCompletedExercises(ctx context.Context, userId, exerciseId string) ([]domain.WorkoutSet, error) {
	res, err := s.workoutSetRepo.GetCompletedExercises(ctx, userId, exerciseId)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed exercises: %w", err)
	}

	return res, nil
}

func (s *trainingService) GetTrainingDays(ctx context.Context, userId string, year, month int) ([]string, error) {
	res, err := s.sessionRepo.GetTrainingDays(ctx, userId, year, month)
	if err != nil {
		return nil, fmt.Errorf("TrainingService: failed to get completed exercises: %w", err)
	}

	return res, nil
}

func (s *trainingService) GetWorkoutsForDay(ctx context.Context, userId, date string) (*types.WorkoutForDay, error) {
	rows, err := s.sessionRepo.GetWorkoutsForDay(ctx, userId, date)
	if err != nil {
		return nil, fmt.Errorf("TrainingService: failed to get workouts for date: %w", err)
	}

	if len(rows) == 0 {
		return &types.WorkoutForDay{
			Date:     date,
			Sessions: []types.WorkoutSession{},
		}, nil
	}

	return mapDayRowsToWorkoutForDay(date, rows), nil
}

func (s *trainingService) GetSessionExercises(ctx context.Context, sessionIDStr string) ([]types.WorkoutExercise, error) {
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid session id format: %w", err)
	}

	rows, err := s.workoutSetRepo.GetSessionExercises(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session exercises: %w", err)
	}

	if len(rows) == 0 {
		return []types.WorkoutExercise{}, nil
	}

	return mapSessionRowsToWorkoutExercises(rows), nil
}

func mapDayRowsToWorkoutForDay(date string, rows []domain.WorkoutDaySetRow) *types.WorkoutForDay {
	type sessionGroup struct {
		session types.WorkoutSession
		exMap   map[uuid.UUID]int // exerciseID -> индекс в slice exercises
	}

	var sessionGroups []sessionGroup
	sessionMap := make(map[uuid.UUID]int) // sessionID -> индекс в slice sessionGroups

	for _, row := range rows {
		sIdx, sExists := sessionMap[row.WorkoutSession.SessionID]
		if !sExists {
			var duration int32
			if row.EndedAt != nil {
				duration = int32(row.EndedAt.Sub(row.StartedAt).Seconds())
			}

			newSession := types.WorkoutSession{
				SessionID:       row.WorkoutSession.SessionID.String(),
				StartedAt:       row.StartedAt,
				EndedAt:         row.EndedAt,
				DurationSeconds: int64(duration),
				Exercises:       []types.WorkoutExercise{},
			}

			sessionGroups = append(sessionGroups, sessionGroup{
				session: newSession,
				exMap:   make(map[uuid.UUID]int),
			})
			sIdx = len(sessionGroups) - 1
			sessionMap[row.WorkoutSession.SessionID] = sIdx
		}

		sGroup := &sessionGroups[sIdx]

		eIdx, eExists := sGroup.exMap[row.ExerciseID]
		if !eExists {
			newExercise := types.WorkoutExercise{
				ExerciseID: row.ExerciseID.String(),
				Name:       row.ExerciseName,
				Type:       string(row.ExerciseType),
				Sets:       []types.CompletedExerciseResponse{},
			}
			sGroup.session.Exercises = append(sGroup.session.Exercises, newExercise)
			eIdx = len(sGroup.session.Exercises) - 1
			sGroup.exMap[row.ExerciseID] = eIdx
		}

		var weight *float32
		if row.Weight != nil {
			w := float32(*row.Weight)
			weight = &w
		}

		var reps *int32
		if row.Reps != nil {
			r := int32(*row.Reps)
			reps = &r
		}

		var durationSec *int32
		if row.DurationSeconds != nil {
			d := int32(*row.DurationSeconds)
			durationSec = &d
		}

		var distanceM *int32
		if row.DistanceMeters != nil {
			d := int32(*row.DistanceMeters)
			distanceM = &d
		}

		set := types.CompletedExerciseResponse{
			SetId:       row.SetID.String(),
			SetNumber:   row.SetNumber,
			Weight:      weight,
			Reps:        reps,
			DurationSec: durationSec,
			DistanceM:   distanceM,
		}

		sGroup.session.Exercises[eIdx].Sets = append(sGroup.session.Exercises[eIdx].Sets, set)
	}

	resultSessions := make([]types.WorkoutSession, len(sessionGroups))
	for i, sg := range sessionGroups {
		resultSessions[i] = sg.session
	}

	return &types.WorkoutForDay{
		Date:     date,
		Sessions: resultSessions,
	}
}

func mapSessionRowsToWorkoutExercises(rows []domain.SessionExerciseRow) []types.WorkoutExercise {
	exMap := make(map[uuid.UUID]int) // exerciseID -> индекс в slice exercises
	exercises := []types.WorkoutExercise{}

	for _, row := range rows {
		eIdx, eExists := exMap[row.ExerciseID]
		if !eExists {
			exercises = append(exercises, types.WorkoutExercise{
				ExerciseID: row.ExerciseID.String(),
				Name:       row.ExerciseName,
				Type:       string(row.ExerciseType),
				Sets:       []types.CompletedExerciseResponse{},
			})
			eIdx = len(exercises) - 1
			exMap[row.ExerciseID] = eIdx
		}

		var weight *float32
		if row.Weight != nil {
			w := float32(*row.Weight)
			weight = &w
		}

		var reps *int32
		if row.Reps != nil {
			r := int32(*row.Reps)
			reps = &r
		}

		var durationSec *int32
		if row.DurationSeconds != nil {
			d := int32(*row.DurationSeconds)
			durationSec = &d
		}

		var distanceM *int32
		if row.DistanceMeters != nil {
			d := int32(*row.DistanceMeters)
			distanceM = &d
		}

		exercises[eIdx].Sets = append(exercises[eIdx].Sets, types.CompletedExerciseResponse{
			SetId:       row.SetID.String(),
			SetNumber:   row.SetNumber,
			Weight:      weight,
			Reps:        reps,
			DurationSec: durationSec,
			DistanceM:   distanceM,
		})
	}

	return exercises
}
