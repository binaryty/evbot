package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

// handleEventDelete ...
func (h *Handler) handleEventDelete(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	defer func() {
		if r := recover(); r != nil {
			h.bot.Send(tgbotapi.NewCallbackWithAlert(query.ID, "⚠️ Произошла ошибка"))
		}
	}()

	chatID := query.Message.Chat.ID

	if !h.isAdmin(query.From.ID) {
		h.sendError(chatID, "🚫 Доступ запрещен")
		return nil
	}

	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		callback := tgbotapi.NewCallbackWithAlert(query.ID, "❌ Ошибка формата запроса")
		h.bot.Send(callback)
	}

	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		callback := tgbotapi.NewCallbackWithAlert(query.ID, "❌ Некорректный ID события")
		h.bot.Send(callback)
		return fmt.Errorf("failed to parse event ID: %w", err)
	}

	// удаляем событие
	if err := h.eventUC.DeleteEvent(ctx, eventID); err != nil {
		callback := tgbotapi.NewCallbackWithAlert(query.ID, "❌ Ошибка удаления события")
		h.bot.Send(callback)
		return err
	}
	h.sendCallback(query.ID, EmOk, "Событие успешно удалено")

	// удаляем сообщение с событием
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID)
	h.bot.Send(deleteMsg)

	return nil
}
