package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleArchiveEvent ...
func (h *Handler) handleArchiveEvent(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	if !h.isAdmin(query.From.ID) {
		h.sendCallback(query.ID, EmCross, "Доступ запрещен")
		return nil
	}

	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		h.sendCallback(query.ID, EmCross, "Некорректный запрос")
		return fmt.Errorf("invalid archive callback payload")
	}

	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse event ID: %w", err)
	}

	h.logger.Debug("archive_event request",
		slog.Int64("adminID", query.From.ID),
		slog.Int64("eventID", eventID))

	if err := h.eventUC.ArchiveEvent(ctx, query.From.ID, eventID); err != nil {
		h.sendCallback(query.ID, EmCross, "Не удалось архивировать событие")
		h.logger.Error("archive event failed",
			slog.Int64("eventID", eventID),
			slog.String("error", err.Error()))
		return fmt.Errorf("archive event failed: %w", err)
	}

	h.logger.Info("event archived",
		slog.Int64("eventID", eventID),
		slog.Int64("adminID", query.From.ID))

	h.sendCallback(query.ID, EmOk, "Событие перенесено в архив")

	deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	if _, err := h.bot.Send(deleteMsg); err != nil {
		h.logger.Warn("failed to delete message after archive", "error", err.Error())
	}

	return nil
}
