package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) handleReportCommand(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	tgIDStr := fmt.Sprintf("%d", msg.From.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	action := tgbotapi.NewChatAction(chatID, tgbotapi.ChatUploadDocument)
	_, _ = r.api.Send(action)

	reportRes, err := r.workoutClient.ExportWorkoutsReport(ctx, tgIDStr, nil, nil)
	if err != nil {
		log.Printf("[ERROR] ExportWorkoutsReport gRPC failed for tgID %s: %v", tgIDStr, err)
		errData := tgbotapi.NewMessage(chatID, "❌ Не удалось сформировать отчет. Убедитесь, что вы зарегистрированы в Mini App.")
		_, _ = r.api.Send(errData)
		return
	}

	filename := reportRes.GetFilename()
	if filename == "" {
		filename = "workouts_report.xlsx"
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: reportRes.GetContent(),
	}

	docMsg := tgbotapi.NewDocument(chatID, fileBytes)
	docMsg.Caption = "📊 Ваш отчет по тренировкам готов!"

	if _, err := r.api.Send(docMsg); err != nil {
		log.Printf("[ERROR] Failed to send document to chat %d: %v", chatID, err)
	}
}
