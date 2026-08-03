package domain

import (
	"time"

	"github.com/google/uuid"
)

// SessionType определяет режим проведения тренировки
type SessionType string

const (
	SessionTypeClassic SessionType = "classic"
	SessionTypeCircuit SessionType = "circuit"
)

// WorkoutSession представляет собой активную или завершенную тренировку
type WorkoutSession struct {
	SessionID  uuid.UUID   `json:"session_id"`
	TenantID   uuid.UUID   `json:"tenant_id"`
	UserID     uuid.UUID   `json:"user_id"`
	TemplateID *uuid.UUID  `json:"template_id,omitempty"`
	Type       SessionType `json:"type"`
	StartedAt  time.Time   `json:"started_at"`
	EndedAt    *time.Time  `json:"ended_at,omitempty"`
	IsActive   bool        `json:"is_active"`
}

type WorkoutDaySetRow struct {
	WorkoutSession `json:"-"`

	ExerciseID   uuid.UUID    `json:"exercise_id"`
	ExerciseName string       `json:"exercise_name"`
	ExerciseType ExerciseType `json:"exercise_type"`

	WorkoutSet `json:"-"`
}

type SessionExerciseRow struct {
	ExerciseID   uuid.UUID    `json:"exercise_id"`
	ExerciseName string       `json:"exercise_name"`
	ExerciseType ExerciseType `json:"exercise_type"`

	WorkoutSet `json:"-"`
}
