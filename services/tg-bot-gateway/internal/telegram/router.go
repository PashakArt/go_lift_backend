package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/clients/workout"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	api           *tgbotapi.BotAPI
	workoutClient *workout.Client
	secretToken   string
	defaultTenant string
}

func NewRouter(api *tgbotapi.BotAPI, workoutClient *workout.Client, secretToken string) *Router {
	return &Router{
		api:           api,
		workoutClient: workoutClient,
		secretToken:   secretToken,
	}
}

func (r *Router) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	if r.secretToken != "" && req.Header.Get("X-Telegram-Bot-Api-Secret-Token") != r.secretToken {
		log.Printf("[WARNING] Unauthorized webhook request from %s", req.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		log.Printf("[ERROR] Failed to decode webhook update: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	go r.processUpdate(update)
}

func (r *Router) processUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			r.handleStartCommand(update.Message)
		case "report":
			r.handleReportCommand(update.Message)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Доступные команды: /start, /report")
			_, _ = r.api.Send(msg)
		}
	}
}

func (r *Router) handleStartCommand(msg *tgbotapi.Message) {
	text := fmt.Sprintf(
		"Привет, %s! 👋\n\n"+
			"Запускай Mini App по кнопке ниже, чтобы начать тренировку, или используй /report для выгрузки Excel-отчета.",
		msg.From.FirstName,
	)

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	_, err := r.api.Send(reply)
	if err != nil {
		log.Printf("[ERROR] Failed to send /start response: %v", err)
	}
}
