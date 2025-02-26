package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

// handleRegistration ...
func (h *Handler) handleRegistration(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	eventID, _ := strconv.ParseInt(strings.Split(query.Data, ":")[1], 10, 64)
	user := domain.User{
		ID:        query.From.ID,
		FirstName: query.From.FirstName,
		UserName:  query.From.UserName,
	}

	isRegistered, err := h.registrationUC.ToggleRegistration(ctx, eventID, &user)
	if err != nil {
		h.sendError(query.Message.Chat.ID, "Ошибка регистрации")
		return fmt.Errorf("failed to register: %w", err)
	}

	isAdmin := h.isAdmin(query.From.ID)

	buttons := createEventButtons(eventID, isRegistered, isAdmin)

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		buttons,
	)
	h.bot.Send(editMarkup)

	return nil
}
