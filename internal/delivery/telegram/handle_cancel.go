package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

// handleCancelCommand ...
func (h *Handler) handleCancelCommand(ctx context.Context, update *tgbotapi.Update) error {
	var userID, chatID int64

	if update.Message != nil && update.Message.From != nil {
		userID = update.Message.From.ID
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
		userID = update.CallbackQuery.From.ID
		if update.CallbackQuery.Message != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}
	} else {
		return errors.New("failed to get user ID")
	}

	if err := h.stateRepo.DeleteState(ctx, userID); err != nil {
		h.sendError(chatID, "Ошибка отмены действия")
		return err
	}

	if update.CallbackQuery != nil {
		parts := strings.Split(update.CallbackQuery.Data, ":")
		eventID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse event ID: %w", err)
		}

		isRegistered, err := h.registrationUC.IsRegistered(ctx, eventID, userID)
		if err != nil {
			return fmt.Errorf("failed to get registraion of user: %w", err)
		}
		isAdmin := h.isAdmin(userID)

		buttons := createEventButtons(eventID, isRegistered, isAdmin)

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			buttons,
		)
		h.bot.Send(editMarkup)
	} else {
		text := "Текущее действие отменено"
		h.sendMsg(chatID, EmOk, text)
	}

	return nil
}
