package types

import (
	"time"

	"github.com/google/uuid"
)

type LogSetRequest struct {
	SessionId  uuid.UUID `json:"session_id"`
	ExerciseId uuid.UUID `json:"exercise_id"`
	SetNumber  int32     `json:"set_number_id"`

	SetID           *uuid.UUID `json:"set_id,omitempty"`
	Weight          *float32   `json:"weight,omitempty"`
	Reps            *int32     `json:"reps,omitempty"`
	DurationSeconds *int32     `json:"duration_seconds,omitempty"`
	DistanceMeters  *int32     `json:"distance_meters,omitempty"`
}

type WorkoutForDay struct {
	Date     string           `json:"date"` // "2026-07-05"
	Sessions []WorkoutSession `json:"sessions"`
}

type WorkoutSession struct {
	SessionID       string            `json:"session_id"`
	StartedAt       time.Time         `json:"started_at"`
	EndedAt         *time.Time        `json:"ended_at,omitempty"`
	DurationSeconds int64             `json:"duration_seconds"`
	Exercises       []WorkoutExercise `json:"exercises"`
}

type WorkoutExercise struct {
	ExerciseID string                      `json:"exercise_id"`
	Name       string                      `json:"name"`
	Type       string                      `json:"type"`
	Sets       []CompletedExerciseResponse `json:"sets"`
}

type CompletedExerciseResponse struct {
	SetNumber   int      `json:"set_number"`
	SetId       string   `json:"set_id"`
	Weight      *float32 `json:"weight,omitempty"`
	Reps        *int32   `json:"reps,omitempty"`
	DurationSec *int32   `json:"duration_sec,omitempty"`
	DistanceM   *int32   `json:"distance_m,omitempty"`
}
