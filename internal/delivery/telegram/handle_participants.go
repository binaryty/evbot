package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"

	"github.com/binaryty/evbot/internal/util"
)

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
		callback := tgbotapi.NewCallbackWithAlert(query.ID, "На событие еще никто не зарегистрирован 🙁")
		h.bot.Send(callback)
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
