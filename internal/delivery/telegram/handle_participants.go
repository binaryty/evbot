package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/binaryty/evbot/internal/util"
)

// handleParticipants ...
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
		return nil
	}

	// Получаем информацию о событии для отображения в сообщении
	event, err := h.eventUC.GetEventByID(ctx, eventID)
	var eventTitle string
	if err == nil && event != nil {
		eventTitle = event.Title
	} else {
		eventTitle = "Событие"
	}

	// Формируем список участников (максимум 10)
	var callbackText strings.Builder
	callbackText.WriteString(fmt.Sprintf("Участники '%s':\n\n", eventTitle))

	maxParticipants := 10
	for i, p := range participants {
		if i >= maxParticipants {
			callbackText.WriteString(fmt.Sprintf("\n... и еще %d человек", len(participants)-maxParticipants))
			break
		}

		name := p.FirstName
		if name == "" && p.UserName != "" {
			name = "@" + p.UserName
		} else if p.UserName != "" {
			name = fmt.Sprintf("%s (@%s)", name, p.UserName)
		}

		callbackText.WriteString(fmt.Sprintf("%d. %s\n", i+1, name))
	}

	// Отправляем колбэк-сообщение с ограниченным списком участников
	callback := tgbotapi.NewCallbackWithAlert(query.ID, callbackText.String())
	_, err = h.bot.Request(callback)
	if err != nil {
		h.logger.Error("failed to send participants callback",
			slog.String("error", err.Error()))
	}

	// Если список слишком большой - предлагаем просмотреть полный список в сообщении
	if len(participants) > maxParticipants {
		// Формируем полный список с экранированием для Markdown
		var fullList strings.Builder
		fullList.WriteString(fmt.Sprintf("👥 *Участники события \"%s\":*\n\n", util.EscapeMarkdownV2(eventTitle)))

		for i, p := range participants {
			// Экранируем спецсимволы
			firstName := util.EscapeMarkdownV2(p.FirstName)
			userName := util.EscapeMarkdownV2(p.UserName)

			var name string
			if firstName != "" && userName != "" {
				name = fmt.Sprintf("%s \\(@%s\\)", firstName, userName)
			} else if firstName != "" {
				name = firstName
			} else if userName != "" {
				name = fmt.Sprintf("@%s", userName)
			} else {
				name = "Участник"
			}

			fullList.WriteString(fmt.Sprintf("%d\\. %s\n", i+1, name))

			// проверяем длину сообщения
			if fullList.Len() > 3000 {
				fullList.WriteString("\n⚠️ Список сокращен из\\-за ограничений Telegram")
				break
			}
		}

		// Отправляем полный список отдельным сообщением
		msg := tgbotapi.NewMessage(chatID, fullList.String())
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		msg.ReplyToMessageID = query.Message.MessageID

		if _, err = h.bot.Send(msg); err != nil {
			h.logger.Error("failed to send full participants list",
				slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}
