package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

func (h *Handler) HandleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update)
	}

	if update.Message == nil {
		return nil
	}

	msg := update.Message
	user := domain.User{
		ID:        msg.From.ID,
		FirstName: msg.From.FirstName,
		UserName:  msg.From.UserName,
	}

	err := h.userUC.CreateOrUpdate(ctx, &user)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
	}

	switch msg.Command() {
	case "start":
		return h.handleStartCommand(ctx, update)
	case "help":
		return h.handleHelpCommand(update)
	case "new_event":
		return h.startNewEvent(ctx, update)
	case "list_events":
		return h.listEvents(ctx, update)
	case "cancel":
		return h.handleCancelCommand(ctx, update)
	default:
		return h.handleUserInput(ctx, update, msg.Text)
	}
}

func GetUserIDFromUpdate(update *tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}

	if update.Message != nil {
		return update.Message.From.ID
	}

	return 0
}
