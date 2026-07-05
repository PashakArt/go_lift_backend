package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutTemplate struct {
	TemplateID  uuid.UUID  `json:"template_id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	CreatorID   *uuid.UUID `json:"creator_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TemplateExercise представляет связь упражнения с шаблоном и его порядковый номер
type TemplateExercise struct {
	TemplateExerciseID uuid.UUID `json:"template_exercise_id"`
	TemplateID         uuid.UUID `json:"template_id"`
	ExerciseID         uuid.UUID `json:"exercise_id"`
	SequenceOrder      int       `json:"sequence_order"`
}
