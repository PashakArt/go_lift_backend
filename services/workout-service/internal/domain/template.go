package domain

import (
	"time"

	"github.com/google/uuid"
)

type TargetSet struct {
	SetNum          int32    `json:"set_num"`
	Weight          *float32 `json:"weight,omitempty"`
	Reps            *int32   `json:"reps,omitempty"`
	DurationSeconds *int32   `json:"duration_seconds,omitempty"`
	DistanceMeters  *int32   `json:"distance_meters,omitempty"`
}

type TemplateItem struct {
	ExerciseID uuid.UUID   `json:"exercise_id"`
	OrderIndex int32       `json:"order_index"`
	TargetSets []TargetSet `json:"target_sets"`
}

type WorkoutTemplate struct {
	TemplateID uuid.UUID      `json:"template_id,omitempty"`
	UserID     uuid.UUID      `json:"user_id"`
	Name       string         `json:"name"`
	Items      []TemplateItem `json:"items"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
}
