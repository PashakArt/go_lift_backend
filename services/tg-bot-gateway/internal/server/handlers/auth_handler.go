package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/types"
)

type AuthHandler struct {
	workoutClient *workout.Client
	jwtManager    *auth.JwtManager
	botToken      string
}

func NewAuthHandler(workoutClient *workout.Client, jwtManager *auth.JwtManager, botToken string) *AuthHandler {
	return &AuthHandler{
		workoutClient: workoutClient,
		jwtManager:    jwtManager,
		botToken:      botToken,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/init", h.HandleInit)
}

func (h *AuthHandler) HandleInit(w http.ResponseWriter, r *http.Request) {
	var req types.InitDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "init data is incorrect")
		return
	}
	defer r.Body.Close()

	var (
		tgIDStr   string
		username  string
		firstName string
		lastName  string
	)

	// 🛠 CHECKS FOR LOCAL DEV MOCK DATA
	// Если мы в dev-режиме и пришли тестовые данные, пропускаем валидацию HMAC
	isDev := os.Getenv("ENV") == "local" || os.Getenv("ENV") == "development"
	if isDev && strings.Contains(req.InitData, "user=%7B%22id%22%3A") {
		log.Println("⚠️ [DEV MODE] Skipping Telegram InitData HMAC validation for mock data")

		tgIDStr = "77777"
		username = "dev_user"
		firstName = "Dev"
		lastName = "Tester"
	} else {
		// PRODUCTION LOGIC
		params, err := ValidateAndParseInitData(req.InitData, h.botToken)
		if err != nil {
			log.Printf("Telegram InitData validation failed: %v\n", err)
			RespondWithError(w, http.StatusUnauthorized, "Telegram InitData validation failed")
			return
		}

		var tgUser types.TelegramUser
		err = json.Unmarshal([]byte(params.Get("user")), &tgUser)
		if err != nil {
			log.Printf("Failed to unmarshal telegram user data: %v\n", err)
			RespondWithError(w, http.StatusBadRequest, "Failed to unmarshal telegram user data")
			return
		}

		tgIDStr = fmt.Sprintf("%d", tgUser.ID)
		username = tgUser.Username
		firstName = tgUser.FirstName
		lastName = tgUser.LastName
	}

	tenantID := req.TenantId
	if tenantID == "" {
		tenantID = defaultTenantId
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.workoutClient.Init(ctx, tenantID, tgIDStr, username, firstName, lastName)
	if err != nil {
		log.Printf("gRPC init failed: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	token, err := h.jwtManager.Generate(res.User.UserId, res.User.TenantId)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	var sessionId string
	var hasActiveSession bool
	var templateId *string
	if res.ActiveSession != nil {
		sessionId = res.ActiveSession.SessionId
		templateId = res.ActiveSession.TemplateId
		hasActiveSession = true
	}

	var statusCode int
	if res.IsNewUser {
		statusCode = http.StatusCreated
	} else {
		statusCode = http.StatusOK
	}

	RespondWithJSON(w, statusCode, types.InitResponse{
		UserID:           res.User.GetUserId(),
		HasActiveSession: hasActiveSession,
		SessionId:        sessionId,
		Token:            token,
		IsNewUser:        res.IsNewUser,
		Branding:         res.TenantBranding,
		TemplateID:       templateId,
		TgUsername:       res.User.TgUsername,
		TgFirstName:      res.User.TgFirstName,
	})
}
