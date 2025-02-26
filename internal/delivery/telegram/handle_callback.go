package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func (h *Handler) handleCallback(ctx context.Context, update *tgbotapi.Update) error {
	query := update.CallbackQuery
	// убрать анимацию кнопки
	defer func() {
		callback := tgbotapi.NewCallback(query.ID, "")
		if _, err := h.bot.Request(callback); err != nil {
			log.Printf("failed to send callback: %v", err)
		}
	}()

	// обработка callback
	data := query.Data
	parts := strings.Split(data, ":")
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "register":
		return h.handleRegistration(ctx, query)
	case "participants":
		return h.handleParticipants(ctx, query)
	case "calendar":
		return h.handleCalendarCallback(ctx, query)
	case "delete_confirm":
		return h.handleDeleteConfirmation(ctx, update)
	case "delete_event":
		return h.handleEventDelete(ctx, query)
	case "delete_cancel":
		return h.handleCancelCommand(ctx, update)
	}

	return nil
}
