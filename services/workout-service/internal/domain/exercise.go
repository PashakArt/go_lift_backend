package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExerciseType string

const (
	ExerciseTypeDynamic    ExerciseType = "dynamic"
	ExerciseTypeStatic     ExerciseType = "static"
	ExerciseTypeBodyweight ExerciseType = "bodyweight"
	ExerciseTypeCardio     ExerciseType = "cardio"
)

type Exercise struct {
	ExerciseID       uuid.UUID    `json:"exercise_id"`
	TenantID         uuid.UUID    `json:"tenant_id"`
	Name             string       `json:"name"`
	Type             ExerciseType `json:"type"`
	MuscleGroupCodes []string     `json:"muscle_group_codes"`
	IsGlobal         bool         `json:"is_global"`
	CreatedAt        time.Time    `json:"created_at"`
}
