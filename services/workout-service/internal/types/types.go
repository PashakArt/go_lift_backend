package types

import (
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
