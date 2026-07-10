package bot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
)

const (
	defaultTenantId = "00000000-0000-0000-0000-000000000000"
)

type InitDataRequest struct {
	InitData string `json:"init_data"`
}

type AuthResponse struct {
	UserID           string `json:"user_id"`
	Token            string `json:"token"`
	HasActiveSession bool   `json:"has_active_session"`
	SessionId        string `json:"session_id"`
}

type HTTPServer struct {
	workoutClient *workout.Client
	botToken      string
}

func NewHTTPServer(botToken string, workoutClient *workout.Client) *HTTPServer {
	return &HTTPServer{
		workoutClient: workoutClient,
		botToken:      botToken,
	}
}

func (s *HTTPServer) Start(port string) error {
	http.HandleFunc("/api/v1/auth", s.handleAuth)

	log.Printf("TG Gateway HTTP Server running on port :%s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *HTTPServer) handleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req InitDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params, err := s.validateAndParseInitData(req.InitData)
	if err != nil {
		log.Printf("Telegram InitData validation failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tgUser struct {
		ID int64 `json:"id"`
	}

	err = json.Unmarshal([]byte(params.Get("user")), &tgUser)
	if err != nil {
		log.Printf("Failed to unmarshal telegram user data: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tgIDStr := fmt.Sprintf("%d", tgUser.ID)
	tenantID := params.Get("start_param")

	if tenantID == "" {
		tenantID = defaultTenantId
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.workoutClient.Auth(ctx, tenantID, tgIDStr)
	if err != nil {
		log.Printf("gRPC Auth failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var sessionId string
	var hasActiveSession bool
	if res.ActiveSession != nil {
		sessionId = res.ActiveSession.SessionId
		hasActiveSession = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		UserID:           res.User.GetUserId(),
		HasActiveSession: hasActiveSession,
		SessionId:        sessionId,
		// TODO: генерировать реальный токен
		Token: "mock-jwt-token",
	})
}

func (s *HTTPServer) validateAndParseInitData(initData string) (url.Values, error) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	// TODO
	if s.botToken == "mock_token_123456" {
		log.Println("[DEBUG] Running in mock mode, skipping Telegram hash validation")
		return params, nil
	}

	hash := params.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash is missing")
	}

	var keys []string
	for k := range params {
		if k != "hash" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var checkStrings []string
	for _, k := range keys {
		checkStrings = append(checkStrings, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}
	checkString := strings.Join(checkStrings, "\n")

	macKey := hmac.New(sha256.New, []byte("WebAppData"))
	macKey.Write([]byte(s.botToken))
	secretKey := macKey.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(checkString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if hash != expectedHash {
		return nil, fmt.Errorf("invalid hash signature")
	}

	return params, nil
}
