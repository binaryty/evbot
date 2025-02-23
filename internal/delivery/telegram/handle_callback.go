package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/util"
)

func (h *Handler) handleCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	// убрать анимацию кнопки
	defer func() {
		callback := tgbotapi.NewCallback(query.ID, "")
		if _, err := h.bot.Request(callback); err != nil {
			log.Printf("failed to send callback: %v", err)
		}
	}()

	// обработка callback
	data := query.Data
	parts := strings.Split(data, ":")
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "register":
		return h.handleRegistration(ctx, query)
	case "participants":
		return h.handleParticipants(ctx, query)
	case "calendar":
		return h.handleCalendarCallback(ctx, query)
	}

	return nil
}

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

	newText := "🎫 Регистрация"
	if isRegistered {
		newText = "✅ Зарегистрирован"
	}

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(newText, fmt.Sprintf("register:%d", eventID)),
				tgbotapi.NewInlineKeyboardButtonData("👥 Участники",
					fmt.Sprintf("participants:%d", eventID)),
			),
		),
	)

	_, err = h.bot.Send(editMarkup)
	return err
}

func (h *Handler) handleParticipants(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	chatID := query.Message.Chat.ID

	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		h.sendError(chatID, "Ошибка обработки события")
		return fmt.Errorf("invalid callback format: %s", query.Data)
	}

	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		h.sendError(chatID, "Ошибка обработки запроса")
		return fmt.Errorf("failed to parse event ID: %w", err)
	}

	// Получаем список участников
	participants, err := h.registrationUC.GetParticipants(ctx, eventID)
	if err != nil {
		h.sendError(chatID, "Ошибка получения участников")
		return fmt.Errorf("failed to get list of participants: %w", err)
	}

	if len(participants) == 0 {
		msg := tgbotapi.NewMessage(chatID, "На событие еще никто не зарегистрирован 🙁")
		_, err := h.bot.Send(msg)
		return err
	}

	// Формируем список с экранированием
	var list strings.Builder
	list.WriteString("👥 *Участники события:*\n\n")

	for _, p := range participants {
		// Экранируем спецсимволы
		firstName := util.EscapeMarkdownV2(p.FirstName)
		userName := util.EscapeMarkdownV2(p.UserName)

		list.WriteString(fmt.Sprintf("• %s \\(@%s\\)\n", firstName, userName))

		// проверяем длину сообщения
		if list.Len() > 3000 {
			list.WriteString("\n⚠️ Список сокращен из-за ограничений Telegram")
			break
		}
	}

	// Отправляем список
	msg := tgbotapi.NewMessage(chatID, list.String())
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.ReplyToMessageID = query.Message.MessageID

	if _, err = h.bot.Send(msg); err != nil {
		log.Printf("Failed to send participants list: %v", err)
		return err
	}

	return nil
}
