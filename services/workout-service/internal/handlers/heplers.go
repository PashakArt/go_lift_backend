package handlers

import (
	workoutv1 "github.com/PashakArt/go_lift_backend/api/proto/workout/v1"
	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
)

func MapDomainTypeToProto(t domain.ExerciseType) workoutv1.ExerciseType {
	switch t {
	case domain.ExerciseTypeDynamic:
		return workoutv1.ExerciseType_EXERCISE_TYPE_DYNAMIC
	case domain.ExerciseTypeStatic:
		return workoutv1.ExerciseType_EXERCISE_TYPE_STATIC
	case domain.ExerciseTypeBodyweight:
		return workoutv1.ExerciseType_EXERCISE_TYPE_BODYWEIGHT
	case domain.ExerciseTypeCardio:
		return workoutv1.ExerciseType_EXERCISE_TYPE_CARDIO
	default:
		return workoutv1.ExerciseType_EXERCISE_TYPE_UNSPECIFIED
	}
}

func MapProtoTypeToDomain(t workoutv1.SessionType) domain.SessionType {
	switch t {
	case workoutv1.SessionType_SESSION_TYPE_CLASSIC:
		return domain.SessionTypeClassic
	case workoutv1.SessionType_SESSION_TYPE_CIRCUIT:
		return domain.SessionTypeCircuit
	default:
		return domain.SessionTypeClassic
	}
}
