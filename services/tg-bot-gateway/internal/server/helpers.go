package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *HTTPServer) ValidateAndParseInitData(initData string) (url.Values, error) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initData query: %w", err)
	}

	if s.botToken == "mock_token_123456" || s.botToken == "" {
		log.Println("[WARN] Running initData validation in MOCK mode")
		return params, nil
	}

	hash := params.Get("hash")
	if hash == "" {
		return nil, errors.New("hash is missing in initData")
	}

	authDateStr := params.Get("auth_date")
	if authDateStr == "" {
		return nil, errors.New("auth_date is missing")
	}

	authTimestamp, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth_date format")
	}

	// Считаем данные просроченными, если им больше 24 часов (86400 секунд)
	if time.Now().Unix()-authTimestamp > 86400 {
		return nil, errors.New("initData is outdated (expired)")
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
	calculatedHashBytes := mac.Sum(nil)

	receivedHashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, errors.New("failed to decode hash from hex")
	}

	if !hmac.Equal(calculatedHashBytes, receivedHashBytes) {
		return nil, errors.New("invalid hash signature")
	}

	return params, nil
}

func RespondWithJSON[T any](w http.ResponseWriter, status int, payload T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("[ERROR] Failed to encode JSON response: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}
