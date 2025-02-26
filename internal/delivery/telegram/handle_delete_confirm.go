package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

func (h *Handler) handleDeleteConfirmation(ctx context.Context, update *tgbotapi.Update) error {
	query := update.CallbackQuery
	chatID := query.Message.Chat.ID

	if !h.isAdmin(query.From.ID) {
		h.sendError(chatID, "Доступ запрещен")
		return nil
	}

	parts := strings.Split(query.Data, ":")
	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Printf("failed to parse event ID: %v", err)
		return fmt.Errorf("failed to parse event ID: %w", err)
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"✅ Подтвердить удаление",
				fmt.Sprintf("delete_event:%d", eventID),
			),

			tgbotapi.NewInlineKeyboardButtonData(
				"❌ Отмена",
				fmt.Sprintf("delete_cancel:%d", eventID),
			),
		),
	)

	edit := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		query.Message.MessageID,
		markup,
	)
	h.bot.Send(edit)

	return nil
}
