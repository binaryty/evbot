package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleMenuButtons обрабатывает нажатия на кнопки главного меню
func (h *Handler) handleMenuButtons(ctx context.Context, update *tgbotapi.Update) error {
	msg := update.Message
	text := msg.Text

	switch text {
	case "📋 Список событий":
		return h.listEvents(ctx, update)
	case "📦 Архив событий":
		return h.listArchivedEvents(ctx, update)
	case "ℹ️ Помощь":
		return h.handleHelpCommand(ctx, update)
	case "🆕 Создать событие":
		return h.startNewEvent(ctx, update)
	default:
		return nil
	}
}
