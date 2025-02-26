package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleHelpCommand(update *tgbotapi.Update) error {
	helpText := `📖 *Справка по командам*

*/new_event* - начать создание нового события
*/list_events* - показать список всех событий с кнопками управления
*/cancel* - отменить текущую операцию
*/help* - показать эту справку

*Как это работает:*
1. Создайте событие с помощью */new_event*
2. В списке событий (*/list_events*) вы можете:
   - 🎫 Зарегистрироваться на событие
   - 👥 Посмотреть список участников
3. Управляйте регистрациями через интерактивные кнопки`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"
	_, err := h.bot.Send(msg)
	return err
}
