package types

import (
	"time"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
)

type InitDataRequest struct {
	InitData string `json:"init_data"`
	TenantId string `json:"tenant_id"`
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

type TrainingDaysResponse struct {
	Days []string `json:"days"`
}

type LogSetRequest struct {
	SessionId  string `json:"session_id" validate:"required,uuid4"`
	ExerciseId string `json:"exercise_id" validate:"required,uuid4"`

	SetID           *string  `json:"set_id,omitempty" validate:"omitempty,uuid4"`
	Weight          *float32 `json:"weight,omitempty" validate:"omitempty,gte=0"`
	Reps            *int32   `json:"reps,omitempty" validate:"omitempty,gte=0"`
	DurationSeconds *int32   `json:"duration_seconds,omitempty" validate:"omitempty,gte=0"`
	DistanceMeters  *int32   `json:"distance_meters,omitempty" validate:"omitempty,gte=0"`
}

type LogSetResponse struct {
	SetID     string `json:"set_id"`
	SetNumber int32  `json:"set_number"`
}

type InitResponse struct {
	UserID           string                    `json:"user_id"`
	Token            string                    `json:"token"`
	HasActiveSession bool                      `json:"has_active_session"`
	IsNewUser        bool                      `json:"is_new_user"`
	SessionId        string                    `json:"session_id"`
	Branding         *workoutv1.TenantBranding `json:"branding"`
	TemplateID       *string                   `json:"template_id,omitempty"`
	TgUsername       *string                   `json:"tg_username"`
	TgFirstName      *string                   `json:"tg_first_name"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type MuscleGroupsResponse struct {
	MuscleGroupId string `json:"muscle_group_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
}

type ExerciseResponse struct {
	ExerciseId string `json:"exercise_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

type StartTrainingResponse struct {
	SessionId string `json:"session_id"`
}

type CompletedExerciseResponse struct {
	SetNumber       int      `json:"set_number"`
	SetId           string   `json:"set_id"`
	Weight          *float32 `json:"weight,omitempty"`
	Reps            *int32   `json:"reps,omitempty"`
	DurationSeconds *int32   `json:"duration_seconds,omitempty"`
	DistanceM       *int32   `json:"distance_m,omitempty"`
}

type CreateTemplateRequest struct {
	Name  string               `json:"name" validate:"required,min=1,max=255"`
	Items []CreateTemplateItem `json:"items" validate:"required,min=1,dive"`
}

type CreateTemplateItem struct {
	ExerciseID string      `json:"exercise_id" validate:"required,uuid4"`
	OrderIndex int32       `json:"order_index" validate:"gte=0"`
	TargetSets []TargetSet `json:"target_sets" validate:"required,min=1,dive"`
}
type TargetSet struct {
	SetNum          int32    `json:"set_num" validate:"required,gt=0"`
	Weight          *float32 `json:"weight,omitempty" validate:"omitempty,gte=0"`
	Reps            *int32   `json:"reps,omitempty" validate:"omitempty,gte=0"`
	DurationSeconds *int32   `json:"duration_seconds,omitempty" validate:"omitempty,gte=0"`
	DistanceMeters  *int32   `json:"distance_meters,omitempty" validate:"omitempty,gte=0"`
}

type CreateTemplateResponse struct {
	TemplateID string `json:"template_id"`
}

type TemplateSummaryResponse struct {
	TemplateID     string `json:"template_id"`
	Name           string `json:"name"`
	ExercisesCount int    `json:"exercises_count"`
	CreatedAt      string `json:"created_at"`
}

type TemplateDetailResponse struct {
	TemplateID string               `json:"template_id"`
	Name       string               `json:"name"`
	CreatedAt  string               `json:"created_at"`
	Items      []TemplateDetailItem `json:"items"`
}

type TemplateDetailItem struct {
	ExerciseID string      `json:"exercise_id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	OrderIndex int32       `json:"order_index"`
	TargetSets []TargetSet `json:"target_sets"`
}

type UpdateTemplateRequest struct {
	Name  string               `json:"name" validate:"required,min=1,max=100"`
	Items []TemplateDetailItem `json:"items" validate:"required,dive"`
}

type StartTrainingRequest struct {
	TemplateID string `json:"template_id,omitempty"`
}

type ReportQueryParams struct {
	FromDate string `json:"from" validate:"omitempty,datetime=2006-01-02"`
	ToDate   string `json:"to"   validate:"omitempty,datetime=2006-01-02"`
}
