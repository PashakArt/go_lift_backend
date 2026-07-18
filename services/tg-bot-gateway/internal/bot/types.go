package bot

import workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"

type InitDataRequest struct {
	InitData string `json:"init_data"`
	TenantId string `json:"tenant_id"`
}

type InitResponse struct {
	UserID           string                    `json:"user_id"`
	Token            string                    `json:"token"`
	HasActiveSession bool                      `json:"has_active_session"`
	IsNewUser        bool                      `json:"is_new_user"`
	SessionId        string                    `json:"session_id"`
	Branding         *workoutv1.TenantBranding `json:"branding"`
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
