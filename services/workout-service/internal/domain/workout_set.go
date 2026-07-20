package domain

import (
	"time"

	"github.com/google/uuid"
)

// WorkoutSet представляет собой один выполненный подход или круг упражнения
type WorkoutSet struct {
	SetID       uuid.UUID `json:"set_id"`
	SessionID   uuid.UUID `json:"session_id"`
	ExerciseID  uuid.UUID `json:"exercise_id"`
	RoundNumber int       `json:"round_number"`
	SetNumber   int       `json:"set_number"`

	// Полиморфные nullable-поля (зависят от ExerciseType)
	Weight          *float64 `json:"weight,omitempty"`           // Для dynamic (в кг, например 82.5)
	Reps            *int     `json:"reps,omitempty"`             // Для dynamic и bodyweight
	DurationSeconds *int     `json:"duration_seconds,omitempty"` // Для static и cardio
	DistanceMeters  *int     `json:"distance_meters,omitempty"`  // Для cardio

	CreatedAt time.Time `json:"created_at"`
}
