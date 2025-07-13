package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCancelCommand ...
func (h *Handler) handleCancelCommand(ctx context.Context, update *tgbotapi.Update) error {
	userID, chatID, err := extractUserAndChatID(*update)
	if err != nil {
		return fmt.Errorf("failed to extract user/chat ID: %w", err)
	}

	// Проверим, было ли активное состояние
	state, getErr := h.stateRepo.GetState(ctx, userID)

	if err := h.stateRepo.DeleteState(ctx, userID); err != nil {
		h.sendError(chatID, "Ошибка отмены действия")
		h.logger.Error("failed to delete state", "userID", userID, "error", err.Error())
		return err
	}

	if update.CallbackQuery != nil {
		if err := h.HandleCancelCallback(ctx, update.CallbackQuery, userID); err != nil {
			h.logger.Warn("failed to handle cancel callback", "error", err.Error())
		}
		return nil
	}

	var text string

	if getErr == nil && state != nil && state.Step != "" {
		// Если было активное состояние, сообщаем какой процесс отменен
		switch state.Step {
		case domain.StepTitle, domain.StepDescription, domain.StepDate, domain.StepTime:
			text = "🚫 Создание события отменено"
		default:
			text = "⚠️ Текущее действие отменено"
		}
	} else {
		text = "⚠️ Нет активных действий для отмены"
	}

	h.sendMsg(chatID, "", text)

	return nil
}

// HandleCancelCallback ...
func (h *Handler) HandleCancelCallback(ctx context.Context, query *tgbotapi.CallbackQuery, userID int64) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid callback data format: %v", query.Data)
	}
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
		query.Message.Chat.ID,
		query.Message.MessageID,
		buttons,
	)

	if _, err := h.bot.Send(editMarkup); err != nil {
		h.logger.Warn("failed to edit markup after cancel", "error", err.Error())
	}

	return nil
}

// extractUserAndChatID ...
func extractUserAndChatID(update tgbotapi.Update) (int64, int64, error) {
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
		return 0, 0, errors.New("failed to get user ID")
	}

	return userID, chatID, nil
}
