package bot

type InitDataRequest struct {
	InitData string `json:"init_data"`
}

type AuthResponse struct {
	UserID           string `json:"user_id"`
	Token            string `json:"token"`
	HasActiveSession bool   `json:"has_active_session"`
	SessionId        string `json:"session_id"`
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
