package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// createMainMenu создает главное меню бота (клавиатуру) в зависимости от прав пользователя
func (h *Handler) createMainMenu(userID int64) tgbotapi.ReplyKeyboardMarkup {
	isAdmin := h.isAdmin(userID)

	var keyboard [][]tgbotapi.KeyboardButton

	// Общие команды для всех пользователей
	keyboard = append(keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("📋 Список событий"),
		tgbotapi.NewKeyboardButton("📦 Архив событий"),
	})
	keyboard = append(keyboard, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
	})

	// Команды только для администраторов
	if isAdmin {
		keyboard = append(keyboard, []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("🆕 Создать событие"),
		})
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

// sendMainMenu отправляет главное меню пользователю
func (h *Handler) sendMainMenu(chatID int64, userID int64) error {
	menu := h.createMainMenu(userID)
	menu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, "Меню загружено")
	msg.ReplyMarkup = menu

	_, err := h.bot.Send(msg)
	return err
}

// hideKeyboard удаляет клавиатуру (если нужно)
func (h *Handler) hideKeyboard(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, "Клавиатура скрыта")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	_, err := h.bot.Send(msg)
	return err
}

// isMenuButton проверяет, является ли текст командой из главного меню
func (h *Handler) isMenuButton(text string) bool {
	switch text {
	case "📋 Список событий", "📦 Архив событий", "ℹ️ Помощь", "🆕 Создать событие":
		return true
	default:
		return false
	}
}

// handleMenuCommand обрабатывает команду /menu
func (h *Handler) handleMenuCommand(update *tgbotapi.Update) error {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Показываем клавиатуру меню
	menu := h.createMainMenu(userID)
	menu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, "Главное меню:")
	msg.ReplyMarkup = menu
	_, err := h.bot.Send(msg)

	return err
}
