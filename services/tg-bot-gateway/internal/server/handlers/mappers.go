package handlers

import (
	"time"

	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/types"
)

func MapGetWorkoutsForDayToHTTP(res *workoutv1.GetWorkoutsForDayResponse) types.WorkoutForDay {
	sessions := make([]types.WorkoutSession, 0, len(res.GetSessions()))

	for _, s := range res.GetSessions() {
		exercises := make([]types.WorkoutExercise, 0, len(s.GetExercises()))

		for _, e := range s.GetExercises() {
			sets := make([]types.CompletedExerciseResponse, 0, len(e.GetSets()))

			for _, st := range e.GetSets() {
				var weight *float32
				if st.Weight != nil {
					w := st.GetWeight()
					weight = &w
				}

				var reps *int32
				if st.Reps != nil {
					r := st.GetReps()
					reps = &r
				}

				var durationSec *int32
				if st.DurationSec != nil {
					d := st.GetDurationSec()
					durationSec = &d
				}

				sets = append(sets, types.CompletedExerciseResponse{
					SetId:           st.GetSetId(),
					SetNumber:       int(st.GetSetNumber()),
					Weight:          weight,
					Reps:            reps,
					DurationSeconds: durationSec,
				})
			}

			exercises = append(exercises, types.WorkoutExercise{
				ExerciseID: e.GetExerciseId(),
				Name:       e.GetName(),
				Type:       e.GetType(),
				Sets:       sets,
			})
		}

		var startedAt time.Time
		if s.GetStartedAt() != "" {
			startedAt, _ = time.Parse(time.RFC3339, s.GetStartedAt())
		}

		var endedAt *time.Time
		if s.GetEndedAt() != "" {
			t, err := time.Parse(time.RFC3339, s.GetEndedAt())
			if err == nil {
				endedAt = &t
			}
		}

		sessions = append(sessions, types.WorkoutSession{
			SessionID:       s.GetSessionId(),
			StartedAt:       startedAt,
			EndedAt:         endedAt,
			DurationSeconds: s.GetDurationSeconds(),
			Exercises:       exercises,
		})
	}

	return types.WorkoutForDay{
		Date:     res.GetDate(),
		Sessions: sessions,
	}
}

func MapGetSessionExercisesToHTTP(res *workoutv1.GetSessionExercisesResponse) []types.WorkoutExercise {
	exercises := make([]types.WorkoutExercise, 0, len(res.GetExercises()))

	for _, e := range res.GetExercises() {
		sets := make([]types.CompletedExerciseResponse, 0, len(e.GetSets()))

		for _, st := range e.GetSets() {
			var weight *float32
			if st.Weight != nil {
				w := st.GetWeight()
				weight = &w
			}

			var reps *int32
			if st.Reps != nil {
				r := st.GetReps()
				reps = &r
			}

			var durationSec *int32
			if st.DurationSec != nil {
				d := st.GetDurationSec()
				durationSec = &d
			}

			var distanceM *int32
			if st.DistanceM != nil {
				d := st.GetDistanceM()
				distanceM = &d
			}

			sets = append(sets, types.CompletedExerciseResponse{
				SetId:           st.GetSetId(),
				SetNumber:       int(st.GetSetNumber()),
				Weight:          weight,
				Reps:            reps,
				DurationSeconds: durationSec,
				DistanceM:       distanceM,
			})
		}

		exercises = append(exercises, types.WorkoutExercise{
			ExerciseID: e.GetExerciseId(),
			Name:       e.GetName(),
			Type:       e.GetType(),
			Sets:       sets,
		})
	}

	return exercises
}
