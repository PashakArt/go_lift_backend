package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/types"
	"github.com/go-playground/validator/v10"
)

type TemplateHandler struct {
	validate *validator.Validate

	workoutClient *workout.Client
}

func NewTemplateHandler(
	validate *validator.Validate,
	workoutClient *workout.Client,
) *TemplateHandler {
	return &TemplateHandler{
		validate:      validate,
		workoutClient: workoutClient,
	}
}

func (h *TemplateHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/templates", authMiddleware(http.HandlerFunc(h.HandleCreateTemplate)))
	mux.Handle("GET /api/v1/templates", authMiddleware(http.HandlerFunc(h.HandleGetTemplates)))
	mux.Handle("GET /api/v1/templates/detail/{templateId}", authMiddleware(http.HandlerFunc(h.HandleGetTemplate)))
	mux.Handle("DELETE /api/v1/templates/{templateId}", authMiddleware(http.HandlerFunc(h.HandleDeleteTemplate)))
	mux.Handle("PUT /api/v1/templates/{templateId}", authMiddleware(http.HandlerFunc(h.HandleUpdateTemplate)))
}

func (h *TemplateHandler) HandleCreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req types.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	res, err := h.workoutClient.CreateTemplate(ctx, req, userId)
	if err != nil {
		log.Printf("[ERROR] CreateTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, types.CreateTemplateResponse{
		TemplateID: res.TemplateId,
	})
}

func (h *TemplateHandler) HandleGetTemplates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	grpcRes, err := h.workoutClient.GetTemplates(ctx, userId)
	if err != nil {
		log.Printf("[ERROR] GetTemplates gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	templates := make([]types.TemplateSummaryResponse, 0, len(grpcRes.GetTemplates()))
	for _, t := range grpcRes.GetTemplates() {
		createdAt := ""
		if t.CreatedAt != nil {
			createdAt = t.CreatedAt.AsTime().Format(time.RFC3339)
		}

		templates = append(templates, types.TemplateSummaryResponse{
			TemplateID:     t.GetTemplateId(),
			Name:           t.GetName(),
			ExercisesCount: int(t.GetExercisesCount()),
			CreatedAt:      createdAt,
		})
	}

	RespondWithJSON(w, http.StatusOK, templates)
}

func (h *TemplateHandler) HandleGetTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.PathValue("templateId")
	if templateId == "" {
		RespondWithError(w, http.StatusBadRequest, "templateId is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	grpcRes, err := h.workoutClient.GetTemplate(ctx, templateId, userId)
	if err != nil {
		log.Printf("[ERROR] GetTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	items := make([]types.TemplateDetailItem, 0, len(grpcRes.GetItems()))
	for _, item := range grpcRes.GetItems() {
		targetSets := make([]types.TargetSet, 0, len(item.GetTargetSets()))
		for _, set := range item.GetTargetSets() {
			targetSets = append(targetSets, types.TargetSet{
				SetNum:          int32(set.GetSetNum()),
				Weight:          set.Weight,
				Reps:            set.Reps,
				DurationSeconds: set.DurationSec,
				DistanceMeters:  set.DistanceM,
			})
		}

		items = append(items, types.TemplateDetailItem{
			ExerciseID: item.GetExerciseId(),
			OrderIndex: item.GetOrderIndex(),
			TargetSets: targetSets,
			Name:       item.Name,
			Type:       item.Type.String(),
		})
	}

	createdAt := ""
	if grpcRes.CreatedAt != nil {
		createdAt = grpcRes.CreatedAt.AsTime().Format(time.RFC3339)
	}

	response := types.TemplateDetailResponse{
		TemplateID: grpcRes.GetTemplateId(),
		Name:       grpcRes.GetName(),
		CreatedAt:  createdAt,
		Items:      items,
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func (h *TemplateHandler) HandleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.PathValue("templateId")
	if templateId == "" {
		RespondWithError(w, http.StatusBadRequest, "templateId is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	err := h.workoutClient.DeleteTemplate(ctx, templateId, userId)
	if err != nil {
		log.Printf("[ERROR] DeleteTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON[any](w, http.StatusOK, nil)
}

func (h *TemplateHandler) HandleUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.PathValue("templateId")
	if templateId == "" {
		RespondWithError(w, http.StatusBadRequest, "templateId is required")
		return
	}

	var req types.UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId := auth.UserIDFromContext(ctx)
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "JWT missing user_id")
		return
	}

	err := h.workoutClient.UpdateTemplate(ctx, templateId, userId, req)
	if err != nil {
		log.Printf("[ERROR] UpdateTemplate gRPC call failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON[any](w, http.StatusOK, nil)
}
