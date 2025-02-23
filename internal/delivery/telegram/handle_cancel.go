package telegram

import (
	"context"
)

func (h *Handler) handleCancelCommand(ctx context.Context, userID int64, chatID int64) error {
	if err := h.stateRepo.DeleteState(ctx, userID); err != nil {
		h.sendError(chatID, "Ошибка отмены действия")
		return err
	}

	icon := "✅"
	text := "Текущее действие отменено"
	h.sendMsg(chatID, icon, text)
	return nil
}
