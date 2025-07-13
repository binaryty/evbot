package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleStartCommand(ctx context.Context, update *tgbotapi.Update) error {
	userID := update.Message.From.ID
	isAdmin := h.isAdmin(userID)
	userName := h.getUserName(ctx, userID)

	var welcomeText string

	if isAdmin {
		welcomeText = fmt.Sprintf(
			`👋 Привет, %s! Я бот для управления событиями.

Я могу помочь вам:
✅ Создавать события с напоминаниями
📋 Показывать список ваших событий
👥 Управлять регистрацией участников

*У вас есть права администратора, поэтому вы можете:*
- Создавать новые события
- Удалять существующие события

Основные команды:
*/new_event* - создать новое событие
*/list_events* - показать все события
*/cancel* - отменить текущее действие
*/help* - показать справку
*/menu* - показать главное меню

Начните с создания первого события!`,
			userName,
		)
	} else {
		welcomeText = fmt.Sprintf(
			`👋 Привет, %s! Я бот для управления событиями.

Я могу помочь вам:
📋 Показывать список всех событий
👥 Управлять регистрацией участников

Основные команды:
*/list_events* - показать все события
*/cancel* - отменить текущее действие
*/help* - показать справку
*/menu* - показать главное меню

Вы можете зарегистрироваться на любое событие из списка!`,
			userName,
		)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeText)
	msg.ParseMode = "Markdown"

	// Отправляем приветственное сообщение
	if _, err := h.bot.Send(msg); err != nil {
		return err
	}

	// Отправляем клавиатуру с главным меню
	menu := h.createMainMenu(userID)
	menu.ResizeKeyboard = true

	menuMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используйте кнопки меню для быстрого доступа к функциям:")
	menuMsg.ReplyMarkup = menu
	_, err := h.bot.Send(menuMsg)

	return err
}

func (h *Handler) getUserName(ctx context.Context, userID int64) string {
	user, err := h.userUC.GetUserByID(ctx, userID)
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
