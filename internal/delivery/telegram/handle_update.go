package telegram

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	domain "github.com/binaryty/evbot/internal/domain/entities"
)

// HandleUpdate ...
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

	// Логирование типа чата
	isGroup := msg.Chat.IsGroup() || msg.Chat.IsSuperGroup()
	h.logger.Debug("handleUpdate",
		slog.Any("user", user),
		slog.Int64("chat_id", msg.Chat.ID),
		slog.Bool("is_group_chat", isGroup),
		slog.String("chat_type", msg.Chat.Type))

	if err := h.userUC.CreateOrUpdate(ctx, &user); err != nil {
		h.logger.Error("failed to create or update user",
			slog.Int64("userID", user.ID),
			slog.String("error", err.Error()),
		)
	}

	switch msg.Command() {
	case "start":
		return h.handleStartCommand(ctx, update)
	case "help":
		return h.handleHelpCommand(ctx, update)
	case "menu":
		return h.handleMenuCommand(update)
	case "new_event":
		return h.startNewEvent(ctx, update)
	case "list_events":
		return h.listEvents(ctx, update)
	case "cancel":
		return h.handleCancelCommand(ctx, update)
	default:
		// Проверяем, не является ли текст сообщения нажатием на кнопку меню
		if h.isMenuButton(msg.Text) {
			return h.handleMenuButtons(ctx, update)
		}
		return h.handleUserInput(ctx, msg.From.ID, msg.Chat.ID, msg.Text)
	}
}

// GetUserIDFromUpdate ...
func GetUserIDFromUpdate(update *tgbotapi.Update) int64 {
	switch {
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.ID
	case update.Message != nil:
		return update.Message.From.ID
	default:
		return 0
	}
}
