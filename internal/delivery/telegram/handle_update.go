package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

func (h *Handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}

	if update.Message == nil {
		return nil
	}

	msg := update.Message
	chatID := msg.Chat.ID
	user := domain.User{
		ID:        msg.From.ID,
		FirstName: msg.From.FirstName,
		UserName:  msg.From.UserName,
	}

	err := h.userRepo.CreateOrUpdate(ctx, &user)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
	}

	switch msg.Command() {
	case "start":
		return h.handleStartCommand(ctx, chatID, user.ID)
	case "help":
		return h.handleHelpCommand(chatID)
	case "new_event":
		return h.startNewEvent(ctx, user.ID, msg.Chat.ID)
	case "list_events":
		return h.listEvents(ctx, user.ID, msg.Chat.ID)
	case "cancel":
		return h.handleCancelCommand(ctx, user.ID, chatID)
	default:
		return h.handleUserInput(ctx, user.ID, msg.Chat.ID, msg.Text)
	}
}
