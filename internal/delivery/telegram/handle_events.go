package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/util"
)

func (h *Handler) startNewEvent(ctx context.Context, userID int64, chatID int64) error {
	initialState := domain.EventState{
		Step: domain.StepTitle,
		TempEvent: domain.Event{
			UserID: userID,
		},
	}

	if err := h.stateRepo.SaveState(ctx, userID, initialState); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, "Введите название события:")
	_, err := h.bot.Send(msg)
	return err
}

func (h *Handler) listEvents(ctx context.Context, userID int64, chatID int64) error {
	const op = "handler.listEvents"
	events, err := h.eventUC.ListEvents(ctx)
	if err != nil {
		h.sendError(chatID, "Ошибка получения событий")
		return fmt.Errorf("%s:list events error: %w", op, err)
	}

	if len(events) == 0 {
		msg := tgbotapi.NewMessage(chatID, "У вас пока нет событий")
		_, err = h.bot.Send(msg)
		return err
	}

	var messages []tgbotapi.Chattable

	for _, event := range events {
		// Проверяем регистрацию пользователя
		isRegistered, err := h.registrationUC.IsRegistered(ctx, event.ID, userID)
		if err != nil {
			log.Printf("failed to check if user is registered: %v", err)
			continue
		}

		// Формируем сообщение для каждого события
		eventOwner, _ := h.userUC.User(ctx, event.UserID)

		text := fmt.Sprintf(
			"📌 %s\n"+
				"📝 %s\n"+
				"⏰ %s\n"+
				"*Автор:* %s",
			util.EscapeMarkdownV2(event.Title),
			util.EscapeMarkdownV2(event.Description),
			event.Date.Format("02\\.01\\.2006 15\\:04"),
			util.EscapeMarkdownV2(eventOwner.UserName),
		)

		// Создаем кнопки для каждого события
		regBtnText, regBtnEmoji := "Регистрация", "🎫"
		if isRegistered {
			regBtnText, regBtnEmoji = "Зарегистрирован", "✅"
		}

		btnRow := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", regBtnEmoji, regBtnText),
				fmt.Sprintf("register:%d", event.ID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"👥 Участники",
				fmt.Sprintf("participants:%d", event.ID),
			),
		)

		// Создаем сообщение с кнопками
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(btnRow)
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		messages = append(messages, msg)
	}

	// Отправляем основное сообщение с инструкцией
	infoMsg := tgbotapi.NewMessage(chatID,
		"📋 *Список ваших событий*\n"+
			"Используйте кнопки под каждым событием для управления:")
	infoMsg.ParseMode = "Markdown"
	messages = append([]tgbotapi.Chattable{infoMsg}, messages...)

	// Отправляем все сообщения
	for _, msg := range messages {
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("%s: failed to sending message: %v", op, err)
		}
	}

	return nil
}
