package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleHelpCommand(ctx context.Context, update *tgbotapi.Update) error {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	isAdmin := h.isAdmin(userID)

	var helpText string

	if isAdmin {
		helpText = `📖 *Справка по командам (Администратор)*

*/new_event* - начать создание нового события (только для администраторов)
*/list_events* - показать список всех событий с кнопками управления
*/cancel* - отменить текущую операцию
*/help* - показать эту справку
*/menu* - показать главное меню

*Как это работает:*
1. Создайте событие с помощью */new_event*
2. В списке событий (*/list_events*) вы можете:
   - 🎫 Зарегистрироваться на событие
   - 👥 Посмотреть список участников
   - 🗑️ Удалить событие (только для администраторов)
3. Управляйте регистрациями через интерактивные кнопки`
	} else {
		helpText = `📖 *Справка по командам*

*/list_events* - показать список всех событий с кнопками управления
*/cancel* - отменить текущую операцию
*/help* - показать эту справку
*/menu* - показать главное меню

*Как это работает:*
1. В списке событий (*/list_events*) вы можете:
   - 🎫 Зарегистрироваться на событие
   - 👥 Посмотреть список участников
2. Управляйте регистрациями через интерактивные кнопки`
	}

	msg := tgbotapi.NewMessage(chatID, helpText)
	msg.ParseMode = "Markdown"

	// Отправляем справку
	if _, err := h.bot.Send(msg); err != nil {
		return err
	}

	// Показываем клавиатуру меню
	menu := h.createMainMenu(userID)
	menu.ResizeKeyboard = true

	menuMsg := tgbotapi.NewMessage(chatID, "Используйте кнопки меню для быстрого доступа к функциям:")
	menuMsg.ReplyMarkup = menu
	_, err := h.bot.Send(menuMsg)

	return err
}
