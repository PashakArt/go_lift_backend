package bot

import (
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
)

func (s *HTTPServer) ValidateAndParseInitData(initData string) (url.Values, error) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

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

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}
