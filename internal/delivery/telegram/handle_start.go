package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleStartCommand(ctx context.Context, chatID int64, userID int64) error {
	welcomeText := fmt.Sprintf(
		`👋 Привет, %s! Я бот для управления событиями.

Я могу помочь вам:
✅ Создавать события с напоминаниями
📋 Показывать список ваших событий
👥 Управлять регистрацией участников

Основные команды:
*/new_event* - создать новое событие
*/list_events* - показать все события
*/cancel* - отменить текущее действие
*/help* - показать справку

Начните с создания первого события!`,
		h.getUserName(ctx, userID),
	)

	msg := tgbotapi.NewMessage(chatID, welcomeText)
	msg.ParseMode = "Markdown"
	_, err := h.bot.Send(msg)

	return err
}

func (h *Handler) getUserName(ctx context.Context, userID int64) string {
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "друг"
	}

	if user.FirstName != "" {
		return user.FirstName
	}

	if user.UserName != "" {
		return "@" + user.UserName
	}

	return "друг"
}
