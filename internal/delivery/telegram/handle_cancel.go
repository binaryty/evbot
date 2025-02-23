package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCancelCommand(ctx context.Context, userID int64, chatID int64) error {
	if err := h.stateRepo.DeleteState(ctx, userID); err != nil {
		h.sendError(chatID, "Ошибка отмены действия")
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "Теукщее действие отменено")
	_, err := h.bot.Send(msg)

	return err
}
